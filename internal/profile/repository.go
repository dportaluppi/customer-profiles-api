package profile

import (
	"context"
	"fmt"
	"github.com/aerospike/aerospike-client-go/v6"
	"github.com/dportaluppi/customer-profiles-api/pkg/profile"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type AerospikeRepository struct {
	client    *aerospike.Client
	namespace string
	set       string
}

func NewAerospikeRepository(client *aerospike.Client, namespace string) *AerospikeRepository {
	return &AerospikeRepository{
		client:    client,
		namespace: namespace,
		set:       "profiles",
	}
}

func (r *AerospikeRepository) Updater(ctx context.Context, profile *profile.Profile) (*profile.Profile, error) {
	var err error
	isNew := profile.ID == ""
	if isNew {
		profile.ID = uuid.New().String()
		profile.CreatedAt = now()
	} else {
		profile.UpdatedAt = now()
	}

	key, err := aerospike.NewKey(r.namespace, r.set, profile.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	bins := profileToBins(profile)

	// policy for write operations
	writePolicy := aerospike.NewWritePolicy(0, 0)
	writePolicy.RecordExistsAction = aerospike.UPDATE
	writePolicy.TotalTimeout = 5000 * time.Millisecond // TODO: make this configurable

	err = r.client.PutBins(writePolicy, key, bins...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = r.addSecondaryIndex(profile, writePolicy)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return profile, nil
}

func now() *time.Time {
	n := time.Now().UTC()
	return &n
}

func (r *AerospikeRepository) GetByID(ctx context.Context, id string) (*profile.Profile, error) {
	key, err := aerospike.NewKey(r.namespace, r.set, id)
	if err != nil {
		return nil, err
	}

	readPolicy := aerospike.NewPolicy()
	readPolicy.TotalTimeout = 5000 * time.Millisecond // TODO: make this configurable

	record, err := r.client.Get(readPolicy, key)
	if err != nil {
		return nil, err
	}

	if record == nil {
		return nil, profile.ErrNotFound
	}

	return recordToProfile(record)
}

func (r *AerospikeRepository) Delete(ctx context.Context, id string) error {
	var err error
	key, err := aerospike.NewKey(r.namespace, r.set, id)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	writePolicy := aerospike.NewWritePolicy(0, 0)
	writePolicy.TotalTimeout = 5000 * time.Millisecond // TODO: make this configurable

	// delete the secondary index
	if err = r.deleteSecondaryIndex(ctx, id, writePolicy); err != nil {
		err = errors.WithStack(err)
		return err
	}

	// delete the profile
	existed, err := r.client.Delete(writePolicy, key)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	if !existed {
		return profile.ErrNotFound
	}

	return nil
}

func (r *AerospikeRepository) GetAll(ctx context.Context, page, limit int) ([]*profile.Profile, int, error) {
	// Get the profile IDs from the index {date}:{id}
	start, end := calculateDateRangeForPage(page, limit)
	profileIDs, err := r.getProfileIDsFromIndex(ctx, start, end)
	if err != nil {
		return nil, 0, err
	}

	// Get the profiles by ID
	profiles := make([]*profile.Profile, 0, limit)
	for _, id := range profileIDs {
		p, err := r.GetByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		profiles = append(profiles, p)
	}

	return profiles, len(profiles), nil
}

// calculateDateRangeForPage returns the start and end times for pagination.
// 'end' is the current time, and 'start' is calculated based on the page number and limit.
func calculateDateRangeForPage(page, limit int) (start, end time.Time) {
	end = time.Now()
	start = end.AddDate(0, 0, -1*limit*page)
	return start, end
}

func (r *AerospikeRepository) getProfileIDsFromIndex(ctx context.Context, start, end time.Time) ([]string, error) {
	var profileIDs []string

	// Define scan policy to include bin data in the results
	scanPolicy := aerospike.NewScanPolicy()
	scanPolicy.IncludeBinData = true

	// Perform the scan on the secondary index set
	recordset, err := r.client.ScanAll(scanPolicy, r.namespace, "profileIndexSet")
	if err != nil {
		return nil, err
	}

	for result := range recordset.Results() {
		if result.Err != nil {
			return nil, result.Err
		}

		// Skip if the record is nil
		if result.Record == nil {
			continue
		}

		// Assume the bins include a reference to the profile ID
		profileID, ok := result.Record.Bins["profileRef"].(string)
		if !ok {
			continue // Or handle as an error if needed
		}

		// Extract createdAt string from the index key and check if it's within the range
		createdAtStr := getCreatedAtStrFromIndexKey(profileID)
		createdAt, err := time.Parse("20060102", createdAtStr)
		if err != nil {
			continue // Or handle the error as needed
		}

		// Check if the createdAt date is within the specified range
		if createdAt.After(start) && createdAt.Before(end) {
			profileIDs = append(profileIDs, profileID)
		}
	}

	return profileIDs, nil
}

// getCreatedAtStrFromIndexKey extracts the date part from the index key
func getCreatedAtStrFromIndexKey(indexKey string) string {
	parts := strings.Split(indexKey, ":")
	if len(parts) >= 2 {
		return parts[0] // Return the date part
	}
	return ""
}

func (r *AerospikeRepository) addSecondaryIndex(profile *profile.Profile, writePolicy *aerospike.WritePolicy) error {
	// build the key for the index
	indexKeyString := fmt.Sprintf("%s:%s", profile.CreatedAt.Format("20060102"), profile.ID)
	indexKey, err := aerospike.NewKey(r.namespace, "profileIndexSet", indexKeyString)
	if err != nil {
		return errors.WithStack(err)
	}

	// insert or update the index reference
	indexBin := aerospike.NewBin("profileRef", profile.ID)
	err = r.client.PutBins(writePolicy, indexKey, indexBin)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *AerospikeRepository) deleteSecondaryIndex(ctx context.Context, id string, writePolicy *aerospike.WritePolicy) error {
	p, err := r.GetByID(ctx, id)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	indexKeyString := fmt.Sprintf("%s:%s", p.CreatedAt.Format("20060102"), id)
	indexKey, err := aerospike.NewKey(r.namespace, "profileIndexSet", indexKeyString)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	if _, err = r.client.Delete(writePolicy, indexKey); err != nil {
		err = errors.WithStack(err)
		return err
	}

	return nil
}

// ------------ Write ------------
func profileToBins(profile *profile.Profile) []*aerospike.Bin {
	var bins []*aerospike.Bin

	// Bins for simple fields
	bins = append(bins, aerospike.NewBin("profileId", profile.ID))
	bins = append(bins, aerospike.NewBin("name", profile.Name))
	bins = append(bins, aerospike.NewBin("gender", profile.Gender))
	bins = append(bins, aerospike.NewBin("location", profile.Location))
	bins = append(bins, aerospike.NewBin("birthday", profile.Birthday.Unix()))

	// Bins for nested structs
	bins = append(bins, contactToBins(profile.Contact)...)
	bins = append(bins, loyaltyToBins(profile.Loyalty)...)

	// Bins for timestamps
	bins = append(bins, timestampToBin("createdAt", profile.CreatedAt))
	bins = append(bins, timestampToBin("updatedAt", profile.UpdatedAt))

	return bins
}

func contactToBins(contact profile.Contact) []*aerospike.Bin {
	return []*aerospike.Bin{
		aerospike.NewBin("contactEmail", contact.Email),
		aerospike.NewBin("contactPhone", contact.Phone),
	}
}

func loyaltyToBins(loyalty profile.Loyalty) []*aerospike.Bin {
	return []*aerospike.Bin{
		aerospike.NewBin("loyaltyLevel", loyalty.Level),
		aerospike.NewBin("enrolledAt", loyalty.EnrolledAt.Unix()),
		aerospike.NewBin("lastActivityAt", loyalty.LastActivityAt.Unix()),
	}
}

func timestampToBin(binName string, timestamp *time.Time) *aerospike.Bin {
	if timestamp != nil {
		return aerospike.NewBin(binName, timestamp.Unix())
	}
	return aerospike.NewBin(binName, nil)
}

// ------------ Read ------------
func recordToProfile(record *aerospike.Record) (*profile.Profile, error) {
	if record == nil {
		return nil, profile.ErrNotFound
	}

	birthday, err := getUnixTimeFromRecord(record.Bins["birthday"])
	if err != nil {
		return nil, err
	}

	profile := &profile.Profile{
		ID:        getStringFromRecord(record.Bins["profileId"]),
		Name:      getStringFromRecord(record.Bins["name"]),
		Gender:    getStringFromRecord(record.Bins["gender"]),
		Birthday:  birthday,
		Location:  getStringFromRecord(record.Bins["location"]),
		Contact:   binsToContact(record.Bins),
		Loyalty:   binsToLoyalty(record.Bins),
		CreatedAt: binToTimestamp(record.Bins["createdAt"]),
		UpdatedAt: binToTimestamp(record.Bins["updatedAt"]),
	}

	return profile, nil
}

func binsToContact(bins aerospike.BinMap) profile.Contact {
	return profile.Contact{
		Email: bins["contactEmail"].(string),
		Phone: bins["contactPhone"].(string),
	}
}

func binsToLoyalty(bins aerospike.BinMap) profile.Loyalty {
	loyalty := profile.Loyalty{
		Level: getStringFromRecord(bins["loyaltyLevel"]),
	}

	loyalty.EnrolledAt, _ = getUnixTimeFromRecord(bins["enrolledAt"])
	loyalty.LastActivityAt, _ = getUnixTimeFromRecord(bins["lastActivityAt"])

	return loyalty
}

func binToTimestamp(bin interface{}) *time.Time {
	switch v := bin.(type) {
	case int64:
		t := time.Unix(v, 0)
		return &t
	case int:
		t := time.Unix(int64(v), 0)
		return &t
	default:
		return nil
	}
}

func getUnixTimeFromRecord(value interface{}) (time.Time, error) {
	switch v := value.(type) {
	case int64:
		return time.Unix(v, 0), nil
	case int:
		return time.Unix(int64(v), 0), nil
	default:
		return time.Time{}, fmt.Errorf("invalid type for Unix time: %T", v)
	}
}

func getStringFromRecord(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}
