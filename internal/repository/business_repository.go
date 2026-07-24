package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"starehian-society-platform/internal/models"
)

type BusinessRepository struct {
	db *sqlx.DB
}

func NewBusinessRepository(db *sqlx.DB) *BusinessRepository {
	return &BusinessRepository{db: db}
}

// Business Listings
func (r *BusinessRepository) Create(ctx context.Context, business *models.BusinessListing) error {
	query := `
		INSERT INTO business_listings (id, user_id, business_name, phone_number, location, website, description, instagram_handle, facebook_handle, logo_url, banner_url, is_verified, is_featured, featured_until, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.db.ExecContext(ctx, query,
		business.ID, business.UserID, business.BusinessName, business.PhoneNumber, business.Location,
		business.Website, business.Description, business.InstagramHandle, business.FacebookHandle,
		business.LogoURL, business.BannerURL, business.IsVerified, business.IsFeatured, business.FeaturedUntil, business.Status,
	)
	return err
}

func (r *BusinessRepository) GetByID(ctx context.Context, id string) (*models.BusinessListing, error) {
	var business models.BusinessListing
	query := `SELECT * FROM business_listings WHERE id = $1`
	err := r.db.GetContext(ctx, &business, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &business, err
}

func (r *BusinessRepository) GetByUserID(ctx context.Context, userID string) ([]*models.BusinessListing, error) {
	var businesses []*models.BusinessListing
	query := `SELECT * FROM business_listings WHERE user_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &businesses, query, userID)
	return businesses, err
}

func (r *BusinessRepository) List(ctx context.Context, limit, offset int) ([]*models.BusinessListing, error) {
	var businesses []*models.BusinessListing
	query := `SELECT * FROM business_listings WHERE status = 'active' ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &businesses, query, limit, offset)
	return businesses, err
}

func (r *BusinessRepository) Update(ctx context.Context, business *models.BusinessListing) error {
	query := `
		UPDATE business_listings 
		SET business_name = $2, phone_number = $3, location = $4, website = $5, description = $6, 
		    instagram_handle = $7, facebook_handle = $8, logo_url = $9, banner_url = $10, 
		    is_verified = $11, is_featured = $12, featured_until = $13, status = $14
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		business.ID, business.BusinessName, business.PhoneNumber, business.Location, business.Website,
		business.Description, business.InstagramHandle, business.FacebookHandle, business.LogoURL, business.BannerURL,
		business.IsVerified, business.IsFeatured, business.FeaturedUntil, business.Status,
	)
	return err
}

func (r *BusinessRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM business_listings WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Jobs
func (r *BusinessRepository) CreateJob(ctx context.Context, job *models.Job) error {
	query := `
		INSERT INTO jobs (id, user_id, business_id, title, description, requirements, responsibilities, location, job_type, salary_range, application_deadline, status, views_count, applications_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.db.ExecContext(ctx, query,
		job.ID, job.UserID, job.BusinessID, job.Title, job.Description, job.Requirements, job.Responsibilities,
		job.Location, job.JobType, job.SalaryRange, job.ApplicationDeadline, job.Status, job.ViewsCount, job.ApplicationsCount,
	)
	return err
}

func (r *BusinessRepository) GetJobByID(ctx context.Context, id string) (*models.Job, error) {
	var job models.Job
	query := `SELECT * FROM jobs WHERE id = $1`
	err := r.db.GetContext(ctx, &job, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &job, err
}

func (r *BusinessRepository) ListJobs(ctx context.Context, limit, offset int) ([]*models.Job, error) {
	var jobs []*models.Job
	query := `SELECT * FROM jobs WHERE status = 'active' ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &jobs, query, limit, offset)
	return jobs, err
}

func (r *BusinessRepository) UpdateJob(ctx context.Context, job *models.Job) error {
	query := `
		UPDATE jobs 
		SET title = $2, description = $3, requirements = $4, responsibilities = $5, location = $6, 
		    job_type = $7, salary_range = $8, application_deadline = $9, status = $10, views_count = $11, applications_count = $12
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		job.ID, job.Title, job.Description, job.Requirements, job.Responsibilities, job.Location,
		job.JobType, job.SalaryRange, job.ApplicationDeadline, job.Status, job.ViewsCount, job.ApplicationsCount,
	)
	return err
}

func (r *BusinessRepository) IncrementJobViews(ctx context.Context, id string) error {
	query := `UPDATE jobs SET views_count = views_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Job Applications
func (r *BusinessRepository) CreateJobApplication(ctx context.Context, application *models.JobApplication) error {
	query := `
		INSERT INTO job_applications (id, job_id, user_id, cover_letter, resume_url, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		application.ID, application.JobID, application.UserID, application.CoverLetter, application.ResumeURL, application.Status,
	)
	return err
}

func (r *BusinessRepository) GetJobApplication(ctx context.Context, jobID, userID string) (*models.JobApplication, error) {
	var application models.JobApplication
	query := `SELECT * FROM job_applications WHERE job_id = $1 AND user_id = $2`
	err := r.db.GetContext(ctx, &application, query, jobID, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &application, err
}

func (r *BusinessRepository) ListJobApplications(ctx context.Context, jobID string) ([]*models.JobApplication, error) {
	var applications []*models.JobApplication
	query := `SELECT * FROM job_applications WHERE job_id = $1 ORDER BY applied_at DESC`
	err := r.db.SelectContext(ctx, &applications, query, jobID)
	return applications, err
}

func (r *BusinessRepository) UpdateJobApplication(ctx context.Context, application *models.JobApplication) error {
	query := `UPDATE job_applications SET status = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, application.ID, application.Status)
	return err
}

// Tenders
func (r *BusinessRepository) CreateTender(ctx context.Context, tender *models.Tender) error {
	query := `
		INSERT INTO tenders (id, user_id, organization_name, title, description, requirements, budget_range, submission_deadline, tender_number, category, status, views_count, bids_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.ExecContext(ctx, query,
		tender.ID, tender.UserID, tender.OrganizationName, tender.Title, tender.Description, tender.Requirements,
		tender.BudgetRange, tender.SubmissionDeadline, tender.TenderNumber, tender.Category, tender.Status, tender.ViewsCount, tender.BidsCount,
	)
	return err
}

func (r *BusinessRepository) GetTenderByID(ctx context.Context, id string) (*models.Tender, error) {
	var tender models.Tender
	query := `SELECT * FROM tenders WHERE id = $1`
	err := r.db.GetContext(ctx, &tender, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &tender, err
}

func (r *BusinessRepository) ListTenders(ctx context.Context, limit, offset int) ([]*models.Tender, error) {
	var tenders []*models.Tender
	query := `SELECT * FROM tenders WHERE status = 'open' ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &tenders, query, limit, offset)
	return tenders, err
}

func (r *BusinessRepository) UpdateTender(ctx context.Context, tender *models.Tender) error {
	query := `
		UPDATE tenders 
		SET organization_name = $2, title = $3, description = $4, requirements = $5, budget_range = $6, 
		    submission_deadline = $7, tender_number = $8, category = $9, status = $10, views_count = $11, bids_count = $12
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		tender.ID, tender.OrganizationName, tender.Title, tender.Description, tender.Requirements, tender.BudgetRange,
		tender.SubmissionDeadline, tender.TenderNumber, tender.Category, tender.Status, tender.ViewsCount, tender.BidsCount,
	)
	return err
}

// Tender Bids
func (r *BusinessRepository) CreateTenderBid(ctx context.Context, bid *models.TenderBid) error {
	query := `
		INSERT INTO tender_bids (id, tender_id, user_id, company_name, proposal_description, bid_amount, proposal_url, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		bid.ID, bid.TenderID, bid.UserID, bid.CompanyName, bid.ProposalDescription, bid.BidAmount, bid.ProposalURL, bid.Status,
	)
	return err
}

func (r *BusinessRepository) GetTenderBid(ctx context.Context, tenderID, userID string) (*models.TenderBid, error) {
	var bid models.TenderBid
	query := `SELECT * FROM tender_bids WHERE tender_id = $1 AND user_id = $2`
	err := r.db.GetContext(ctx, &bid, query, tenderID, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &bid, err
}

func (r *BusinessRepository) ListTenderBids(ctx context.Context, tenderID string) ([]*models.TenderBid, error) {
	var bids []*models.TenderBid
	query := `SELECT * FROM tender_bids WHERE tender_id = $1 ORDER BY submitted_at DESC`
	err := r.db.SelectContext(ctx, &bids, query, tenderID)
	return bids, err
}

func (r *BusinessRepository) UpdateTenderBid(ctx context.Context, bid *models.TenderBid) error {
	query := `UPDATE tender_bids SET status = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, bid.ID, bid.Status)
	return err
}

// Class Groups
func (r *BusinessRepository) CreateClassGroup(ctx context.Context, group *models.ClassGroup) error {
	query := `
		INSERT INTO class_groups (id, school_type, year_of_completion, class_name, description, class_rep_id, member_count, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		group.ID, group.SchoolType, group.YearOfCompletion, group.ClassName, group.Description, group.ClassRepID, group.MemberCount, group.IsActive,
	)
	return err
}

func (r *BusinessRepository) GetClassGroupByID(ctx context.Context, id string) (*models.ClassGroup, error) {
	var group models.ClassGroup
	query := `SELECT * FROM class_groups WHERE id = $1`
	err := r.db.GetContext(ctx, &group, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &group, err
}

func (r *BusinessRepository) ListClassGroups(ctx context.Context, schoolType string, year int) ([]*models.ClassGroup, error) {
	var groups []*models.ClassGroup
	query := `SELECT * FROM class_groups WHERE school_type = $1 AND year_of_completion = $2 AND is_active = true ORDER BY class_name`
	err := r.db.SelectContext(ctx, &groups, query, schoolType, year)
	return groups, err
}

func (r *BusinessRepository) UpdateClassGroup(ctx context.Context, group *models.ClassGroup) error {
	query := `
		UPDATE class_groups 
		SET description = $2, class_rep_id = $3, member_count = $4, is_active = $5
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, group.ID, group.Description, group.ClassRepID, group.MemberCount, group.IsActive)
	return err
}

// Class Group Members
func (r *BusinessRepository) AddClassGroupMember(ctx context.Context, member *models.ClassGroupMember) error {
	query := `
		INSERT INTO class_group_members (id, class_group_id, user_id, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (class_group_id, user_id) DO UPDATE SET role = $4
	`
	_, err := r.db.ExecContext(ctx, query, member.ID, member.ClassGroupID, member.UserID, member.Role)
	return err
}

func (r *BusinessRepository) GetClassGroupMember(ctx context.Context, groupID, userID string) (*models.ClassGroupMember, error) {
	var member models.ClassGroupMember
	query := `SELECT * FROM class_group_members WHERE class_group_id = $1 AND user_id = $2`
	err := r.db.GetContext(ctx, &member, query, groupID, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &member, err
}

func (r *BusinessRepository) ListClassGroupMembers(ctx context.Context, groupID string) ([]*models.ClassGroupMember, error) {
	var members []*models.ClassGroupMember
	query := `SELECT * FROM class_group_members WHERE class_group_id = $1 ORDER BY joined_at`
	err := r.db.SelectContext(ctx, &members, query, groupID)
	return members, err
}

func (r *BusinessRepository) RemoveClassGroupMember(ctx context.Context, groupID, userID string) error {
	query := `DELETE FROM class_group_members WHERE class_group_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, groupID, userID)
	return err
}

func (r *BusinessRepository) IncrementClassGroupMemberCount(ctx context.Context, groupID string) error {
	query := `UPDATE class_groups SET member_count = member_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, groupID)
	return err
}

// Merchant Offers
func (r *BusinessRepository) CreateMerchantOffer(ctx context.Context, offer *models.MerchantOffer) error {
	query := `
		INSERT INTO merchant_offers (id, user_id, business_id, title, description, discount_percentage, original_price, offer_price, valid_from, valid_until, terms_conditions, image_url, is_exclusive, status, views_count, redemptions_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	_, err := r.db.ExecContext(ctx, query,
		offer.ID, offer.UserID, offer.BusinessID, offer.Title, offer.Description, offer.DiscountPercentage,
		offer.OriginalPrice, offer.OfferPrice, offer.ValidFrom, offer.ValidUntil, offer.TermsConditions,
		offer.ImageURL, offer.IsExclusive, offer.Status, offer.ViewsCount, offer.RedemptionsCount,
	)
	return err
}

func (r *BusinessRepository) GetMerchantOfferByID(ctx context.Context, id string) (*models.MerchantOffer, error) {
	var offer models.MerchantOffer
	query := `SELECT * FROM merchant_offers WHERE id = $1`
	err := r.db.GetContext(ctx, &offer, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &offer, err
}

func (r *BusinessRepository) ListMerchantOffers(ctx context.Context, limit, offset int) ([]*models.MerchantOffer, error) {
	var offers []*models.MerchantOffer
	query := `SELECT * FROM merchant_offers WHERE status = 'active' AND valid_until > NOW() ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &offers, query, limit, offset)
	return offers, err
}

func (r *BusinessRepository) UpdateMerchantOffer(ctx context.Context, offer *models.MerchantOffer) error {
	query := `
		UPDATE merchant_offers 
		SET title = $2, description = $3, discount_percentage = $4, original_price = $5, offer_price = $6, 
		    valid_from = $7, valid_until = $8, terms_conditions = $9, image_url = $10, is_exclusive = $11, status = $12, 
		    views_count = $13, redemptions_count = $14
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		offer.ID, offer.Title, offer.Description, offer.DiscountPercentage, offer.OriginalPrice, offer.OfferPrice,
		offer.ValidFrom, offer.ValidUntil, offer.TermsConditions, offer.ImageURL, offer.IsExclusive, offer.Status,
		offer.ViewsCount, offer.RedemptionsCount,
	)
	return err
}

// Sponsorships
func (r *BusinessRepository) CreateSponsorship(ctx context.Context, sponsorship *models.Sponsorship) error {
	query := `
		INSERT INTO sponsorships (id, user_id, title, description, sponsorship_type, target_amount, current_amount, start_date, end_date, status, beneficiary, image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		sponsorship.ID, sponsorship.UserID, sponsorship.Title, sponsorship.Description, sponsorship.SponsorshipType,
		sponsorship.TargetAmount, sponsorship.CurrentAmount, sponsorship.StartDate, sponsorship.EndDate, sponsorship.Status,
		sponsorship.Beneficiary, sponsorship.ImageURL,
	)
	return err
}

func (r *BusinessRepository) GetSponsorshipByID(ctx context.Context, id string) (*models.Sponsorship, error) {
	var sponsorship models.Sponsorship
	query := `SELECT * FROM sponsorships WHERE id = $1`
	err := r.db.GetContext(ctx, &sponsorship, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &sponsorship, err
}

func (r *BusinessRepository) ListSponsorships(ctx context.Context, limit, offset int) ([]*models.Sponsorship, error) {
	var sponsorships []*models.Sponsorship
	query := `SELECT * FROM sponsorships WHERE status = 'active' ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &sponsorships, query, limit, offset)
	return sponsorships, err
}

func (r *BusinessRepository) UpdateSponsorship(ctx context.Context, sponsorship *models.Sponsorship) error {
	query := `
		UPDATE sponsorships 
		SET title = $2, description = $3, sponsorship_type = $4, target_amount = $5, current_amount = $6, 
		    start_date = $7, end_date = $8, status = $9, beneficiary = $10, image_url = $11
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		sponsorship.ID, sponsorship.Title, sponsorship.Description, sponsorship.SponsorshipType,
		sponsorship.TargetAmount, sponsorship.CurrentAmount, sponsorship.StartDate, sponsorship.EndDate,
		sponsorship.Status, sponsorship.Beneficiary, sponsorship.ImageURL,
	)
	return err
}

func (r *BusinessRepository) UpdateSponsorshipAmount(ctx context.Context, id string, amount float64) error {
	query := `UPDATE sponsorships SET current_amount = current_amount + $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, amount, id)
	return err
}

// Sponsorship Contributions
func (r *BusinessRepository) CreateSponsorshipContribution(ctx context.Context, contribution *models.SponsorshipContribution) error {
	query := `
		INSERT INTO sponsorship_contributions (id, sponsorship_id, user_id, amount, message, is_anonymous)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		contribution.ID, contribution.SponsorshipID, contribution.UserID, contribution.Amount, contribution.Message, contribution.IsAnonymous,
	)
	return err
}

func (r *BusinessRepository) ListSponsorshipContributions(ctx context.Context, sponsorshipID string) ([]*models.SponsorshipContribution, error) {
	var contributions []*models.SponsorshipContribution
	query := `SELECT * FROM sponsorship_contributions WHERE sponsorship_id = $1 ORDER BY contributed_at DESC`
	err := r.db.SelectContext(ctx, &contributions, query, sponsorshipID)
	return contributions, err
}

// Escrow Transactions
func (r *BusinessRepository) CreateEscrowTransaction(ctx context.Context, transaction *models.EscrowTransaction) error {
	query := `
		INSERT INTO escrow_transactions (id, business_id, buyer_id, seller_id, amount, description, status, release_conditions)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		transaction.ID, transaction.BusinessID, transaction.BuyerID, transaction.SellerID, transaction.Amount,
		transaction.Description, transaction.Status, transaction.ReleaseConditions,
	)
	return err
}

func (r *BusinessRepository) GetEscrowTransactionByID(ctx context.Context, id string) (*models.EscrowTransaction, error) {
	var transaction models.EscrowTransaction
	query := `SELECT * FROM escrow_transactions WHERE id = $1`
	err := r.db.GetContext(ctx, &transaction, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &transaction, err
}

func (r *BusinessRepository) ListEscrowTransactions(ctx context.Context, userID string, role string) ([]*models.EscrowTransaction, error) {
	var transactions []*models.EscrowTransaction
	var query string
	if role == "buyer" {
		query = `SELECT * FROM escrow_transactions WHERE buyer_id = $1 ORDER BY created_at DESC`
	} else {
		query = `SELECT * FROM escrow_transactions WHERE seller_id = $1 ORDER BY created_at DESC`
	}
	err := r.db.SelectContext(ctx, &transactions, query, userID)
	return transactions, err
}

func (r *BusinessRepository) UpdateEscrowTransaction(ctx context.Context, transaction *models.EscrowTransaction) error {
	query := `
		UPDATE escrow_transactions 
		SET status = $2, funded_at = $3, completed_at = $4, cancelled_at = $5
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		transaction.ID, transaction.Status, transaction.FundedAt, transaction.CompletedAt, transaction.CancelledAt,
	)
	return err
}
