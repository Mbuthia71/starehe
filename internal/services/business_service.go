package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"starehian-society-platform/internal/models"
	"starehian-society-platform/internal/repository"
	"starehian-society-platform/pkg/logger"
)

type BusinessService struct {
	businessRepo *repository.BusinessRepository
	logger       *logger.Logger
}

func NewBusinessService(businessRepo *repository.BusinessRepository, logger *logger.Logger) *BusinessService {
	return &BusinessService{
		businessRepo: businessRepo,
		logger:       logger,
	}
}

func getJobType(jobType *string) string {
	if jobType != nil {
		return *jobType
	}
	return "full-time"
}

// Business Listings
func (s *BusinessService) CreateBusinessListing(ctx context.Context, userID, businessName, phoneNumber, location, description string, website, instagram, facebook *string) (*models.BusinessListing, error) {
	business := &models.BusinessListing{
		ID:          uuid.New().String(),
		UserID:      userID,
		BusinessName: businessName,
		PhoneNumber: phoneNumber,
		Location:    location,
		Website:     website,
		Description: description,
		InstagramHandle: instagram,
		FacebookHandle:  facebook,
		Status:      "active",
	}

	err := s.businessRepo.Create(ctx, business)
	if err != nil {
		s.logger.Errorf("Failed to create business listing: %v", err)
		return nil, fmt.Errorf("failed to create business listing: %w", err)
	}

	s.logger.Infof("Business listing created: %s", business.ID)
	return business, nil
}

func (s *BusinessService) GetBusinessListing(ctx context.Context, id string) (*models.BusinessListing, error) {
	return s.businessRepo.GetByID(ctx, id)
}

func (s *BusinessService) GetUserBusinessListings(ctx context.Context, userID string) ([]*models.BusinessListing, error) {
	return s.businessRepo.GetByUserID(ctx, userID)
}

func (s *BusinessService) ListBusinessListings(ctx context.Context, limit, offset int) ([]*models.BusinessListing, error) {
	return s.businessRepo.List(ctx, limit, offset)
}

func (s *BusinessService) UpdateBusinessListing(ctx context.Context, id string, updates map[string]interface{}) error {
	business, err := s.businessRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if business == nil {
		return fmt.Errorf("business listing not found")
	}

	// Apply updates
	if v, ok := updates["business_name"]; ok {
		business.BusinessName = v.(string)
	}
	if v, ok := updates["phone_number"]; ok {
		business.PhoneNumber = v.(string)
	}
	if v, ok := updates["location"]; ok {
		business.Location = v.(string)
	}
	if v, ok := updates["website"]; ok {
		business.Website = v.(*string)
	}
	if v, ok := updates["description"]; ok {
		business.Description = v.(string)
	}
	if v, ok := updates["instagram_handle"]; ok {
		business.InstagramHandle = v.(*string)
	}
	if v, ok := updates["facebook_handle"]; ok {
		business.FacebookHandle = v.(*string)
	}
	if v, ok := updates["logo_url"]; ok {
		business.LogoURL = v.(*string)
	}
	if v, ok := updates["banner_url"]; ok {
		business.BannerURL = v.(*string)
	}
	if v, ok := updates["status"]; ok {
		business.Status = v.(string)
	}

	return s.businessRepo.Update(ctx, business)
}

func (s *BusinessService) DeleteBusinessListing(ctx context.Context, id string) error {
	return s.businessRepo.Delete(ctx, id)
}

// Jobs
func (s *BusinessService) CreateJob(ctx context.Context, userID, title, description string, businessID *string, requirements, responsibilities, location, jobType, salaryRange *string, applicationDeadline *time.Time) (*models.Job, error) {
	job := &models.Job{
		ID:                  uuid.New().String(),
		UserID:              userID,
		BusinessID:          businessID,
		Title:               title,
		Description:         description,
		Requirements:        requirements,
		Responsibilities:    responsibilities,
		Location:            location,
		JobType:             getJobType(jobType),
		SalaryRange:         salaryRange,
		ApplicationDeadline: applicationDeadline,
		Status:              "active",
		ViewsCount:          0,
		ApplicationsCount:   0,
	}

	err := s.businessRepo.CreateJob(ctx, job)
	if err != nil {
		s.logger.Errorf("Failed to create job: %v", err)
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	s.logger.Infof("Job created: %s", job.ID)
	return job, nil
}

func (s *BusinessService) GetJob(ctx context.Context, id string) (*models.Job, error) {
	job, err := s.businessRepo.GetJobByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if job != nil {
		// Increment view count
		_ = s.businessRepo.IncrementJobViews(ctx, id)
	}
	return job, nil
}

func (s *BusinessService) ListJobs(ctx context.Context, limit, offset int) ([]*models.Job, error) {
	return s.businessRepo.ListJobs(ctx, limit, offset)
}

func (s *BusinessService) UpdateJob(ctx context.Context, id string, updates map[string]interface{}) error {
	job, err := s.businessRepo.GetJobByID(ctx, id)
	if err != nil {
		return err
	}
	if job == nil {
		return fmt.Errorf("job not found")
	}

	// Apply updates
	if v, ok := updates["title"]; ok {
		job.Title = v.(string)
	}
	if v, ok := updates["description"]; ok {
		job.Description = v.(string)
	}
	if v, ok := updates["requirements"]; ok {
		job.Requirements = v.(*string)
	}
	if v, ok := updates["responsibilities"]; ok {
		job.Responsibilities = v.(*string)
	}
	if v, ok := updates["location"]; ok {
		job.Location = v.(*string)
	}
	if v, ok := updates["job_type"]; ok {
		job.JobType = v.(string)
	}
	if v, ok := updates["salary_range"]; ok {
		job.SalaryRange = v.(*string)
	}
	if v, ok := updates["application_deadline"]; ok {
		job.ApplicationDeadline = v.(*time.Time)
	}
	if v, ok := updates["status"]; ok {
		job.Status = v.(string)
	}

	return s.businessRepo.UpdateJob(ctx, job)
}

func (s *BusinessService) ApplyForJob(ctx context.Context, jobID, userID string, coverLetter, resumeURL *string) (*models.JobApplication, error) {
	// Check if already applied
	existing, _ := s.businessRepo.GetJobApplication(ctx, jobID, userID)
	if existing != nil {
		return nil, fmt.Errorf("already applied for this job")
	}

	application := &models.JobApplication{
		ID:          uuid.New().String(),
		JobID:       jobID,
		UserID:      userID,
		CoverLetter: coverLetter,
		ResumeURL:   resumeURL,
		Status:      "pending",
	}

	err := s.businessRepo.CreateJobApplication(ctx, application)
	if err != nil {
		s.logger.Errorf("Failed to create job application: %v", err)
		return nil, fmt.Errorf("failed to create job application: %w", err)
	}

	// Increment application count
	job, _ := s.businessRepo.GetJobByID(ctx, jobID)
	if job != nil {
		job.ApplicationsCount++
		_ = s.businessRepo.UpdateJob(ctx, job)
	}

	s.logger.Infof("Job application created: %s", application.ID)
	return application, nil
}

func (s *BusinessService) GetJobApplications(ctx context.Context, jobID string) ([]*models.JobApplication, error) {
	return s.businessRepo.ListJobApplications(ctx, jobID)
}

func (s *BusinessService) UpdateJobApplicationStatus(ctx context.Context, id, status string) error {
	application, err := s.businessRepo.GetJobApplication(ctx, id, "")
	if err != nil {
		return err
	}
	if application == nil {
		return fmt.Errorf("application not found")
	}

	application.Status = status
	return s.businessRepo.UpdateJobApplication(ctx, application)
}

// Tenders
func (s *BusinessService) CreateTender(ctx context.Context, userID, organizationName, title, description string, requirements, budgetRange, tenderNumber, category *string, submissionDeadline time.Time) (*models.Tender, error) {
	tender := &models.Tender{
		ID:                  uuid.New().String(),
		UserID:              userID,
		OrganizationName:    organizationName,
		Title:               title,
		Description:         description,
		Requirements:        requirements,
		BudgetRange:         budgetRange,
		SubmissionDeadline:  submissionDeadline,
		TenderNumber:        tenderNumber,
		Category:            category,
		Status:              "open",
		ViewsCount:          0,
		BidsCount:           0,
	}

	err := s.businessRepo.CreateTender(ctx, tender)
	if err != nil {
		s.logger.Errorf("Failed to create tender: %v", err)
		return nil, fmt.Errorf("failed to create tender: %w", err)
	}

	s.logger.Infof("Tender created: %s", tender.ID)
	return tender, nil
}

func (s *BusinessService) GetTender(ctx context.Context, id string) (*models.Tender, error) {
	return s.businessRepo.GetTenderByID(ctx, id)
}

func (s *BusinessService) ListTenders(ctx context.Context, limit, offset int) ([]*models.Tender, error) {
	return s.businessRepo.ListTenders(ctx, limit, offset)
}

func (s *BusinessService) UpdateTender(ctx context.Context, id string, updates map[string]interface{}) error {
	tender, err := s.businessRepo.GetTenderByID(ctx, id)
	if err != nil {
		return err
	}
	if tender == nil {
		return fmt.Errorf("tender not found")
	}

	// Apply updates
	if v, ok := updates["organization_name"]; ok {
		tender.OrganizationName = v.(string)
	}
	if v, ok := updates["title"]; ok {
		tender.Title = v.(string)
	}
	if v, ok := updates["description"]; ok {
		tender.Description = v.(string)
	}
	if v, ok := updates["requirements"]; ok {
		tender.Requirements = v.(*string)
	}
	if v, ok := updates["budget_range"]; ok {
		tender.BudgetRange = v.(*string)
	}
	if v, ok := updates["submission_deadline"]; ok {
		tender.SubmissionDeadline = v.(time.Time)
	}
	if v, ok := updates["tender_number"]; ok {
		tender.TenderNumber = v.(*string)
	}
	if v, ok := updates["category"]; ok {
		tender.Category = v.(*string)
	}
	if v, ok := updates["status"]; ok {
		tender.Status = v.(string)
	}

	return s.businessRepo.UpdateTender(ctx, tender)
}

func (s *BusinessService) SubmitTenderBid(ctx context.Context, tenderID, userID, companyName, proposalDescription string, bidAmount *float64, proposalURL *string) (*models.TenderBid, error) {
	// Check if already bid
	existing, _ := s.businessRepo.GetTenderBid(ctx, tenderID, userID)
	if existing != nil {
		return nil, fmt.Errorf("already submitted a bid for this tender")
	}

	bid := &models.TenderBid{
		ID:                  uuid.New().String(),
		TenderID:            tenderID,
		UserID:              userID,
		CompanyName:         companyName,
		ProposalDescription: proposalDescription,
		BidAmount:           bidAmount,
		ProposalURL:         proposalURL,
		Status:              "submitted",
	}

	err := s.businessRepo.CreateTenderBid(ctx, bid)
	if err != nil {
		s.logger.Errorf("Failed to create tender bid: %v", err)
		return nil, fmt.Errorf("failed to create tender bid: %w", err)
	}

	// Increment bid count
	tender, _ := s.businessRepo.GetTenderByID(ctx, tenderID)
	if tender != nil {
		tender.BidsCount++
		_ = s.businessRepo.UpdateTender(ctx, tender)
	}

	s.logger.Infof("Tender bid created: %s", bid.ID)
	return bid, nil
}

func (s *BusinessService) GetTenderBids(ctx context.Context, tenderID string) ([]*models.TenderBid, error) {
	return s.businessRepo.ListTenderBids(ctx, tenderID)
}

func (s *BusinessService) UpdateTenderBidStatus(ctx context.Context, id, status string) error {
	bid, err := s.businessRepo.GetTenderBid(ctx, id, "")
	if err != nil {
		return err
	}
	if bid == nil {
		return fmt.Errorf("bid not found")
	}

	bid.Status = status
	return s.businessRepo.UpdateTenderBid(ctx, bid)
}

// Class Groups
func (s *BusinessService) CreateClassGroup(ctx context.Context, schoolType string, yearOfCompletion int, className, description string, classRepID *string) (*models.ClassGroup, error) {
	group := &models.ClassGroup{
		ID:              uuid.New().String(),
		SchoolType:      schoolType,
		YearOfCompletion: yearOfCompletion,
		ClassName:       className,
		Description:     &description,
		ClassRepID:      classRepID,
		MemberCount:     0,
		IsActive:        true,
	}

	err := s.businessRepo.CreateClassGroup(ctx, group)
	if err != nil {
		s.logger.Errorf("Failed to create class group: %v", err)
		return nil, fmt.Errorf("failed to create class group: %w", err)
	}

	s.logger.Infof("Class group created: %s", group.ID)
	return group, nil
}

func (s *BusinessService) GetClassGroup(ctx context.Context, id string) (*models.ClassGroup, error) {
	return s.businessRepo.GetClassGroupByID(ctx, id)
}

func (s *BusinessService) ListClassGroups(ctx context.Context, schoolType string, year int) ([]*models.ClassGroup, error) {
	return s.businessRepo.ListClassGroups(ctx, schoolType, year)
}

func (s *BusinessService) JoinClassGroup(ctx context.Context, groupID, userID string, role string) error {
	// Check if already a member
	existing, _ := s.businessRepo.GetClassGroupMember(ctx, groupID, userID)
	if existing != nil {
		return fmt.Errorf("already a member of this group")
	}

	member := &models.ClassGroupMember{
		ID:           uuid.New().String(),
		ClassGroupID: groupID,
		UserID:       userID,
		Role:         role,
	}

	err := s.businessRepo.AddClassGroupMember(ctx, member)
	if err != nil {
		s.logger.Errorf("Failed to add class group member: %v", err)
		return fmt.Errorf("failed to add class group member: %w", err)
	}

	// Increment member count
	_ = s.businessRepo.IncrementClassGroupMemberCount(ctx, groupID)

	s.logger.Infof("User %s joined class group %s", userID, groupID)
	return nil
}

func (s *BusinessService) GetClassGroupMembers(ctx context.Context, groupID string) ([]*models.ClassGroupMember, error) {
	return s.businessRepo.ListClassGroupMembers(ctx, groupID)
}

func (s *BusinessService) LeaveClassGroup(ctx context.Context, groupID, userID string) error {
	err := s.businessRepo.RemoveClassGroupMember(ctx, groupID, userID)
	if err != nil {
		s.logger.Errorf("Failed to remove class group member: %v", err)
		return fmt.Errorf("failed to remove class group member: %w", err)
	}

	// Decrement member count
	group, _ := s.businessRepo.GetClassGroupByID(ctx, groupID)
	if group != nil && group.MemberCount > 0 {
		group.MemberCount--
		_ = s.businessRepo.UpdateClassGroup(ctx, group)
	}

	s.logger.Infof("User %s left class group %s", userID, groupID)
	return nil
}

// Merchant Offers
func (s *BusinessService) CreateMerchantOffer(ctx context.Context, userID, title, description string, businessID *string, discountPercentage *int, originalPrice, offerPrice *float64, validFrom, validUntil time.Time, termsConditions, imageURL *string, isExclusive bool) (*models.MerchantOffer, error) {
	offer := &models.MerchantOffer{
		ID:                uuid.New().String(),
		UserID:            userID,
		BusinessID:        businessID,
		Title:             title,
		Description:       description,
		DiscountPercentage: discountPercentage,
		OriginalPrice:     originalPrice,
		OfferPrice:        offerPrice,
		ValidFrom:         validFrom,
		ValidUntil:        validUntil,
		TermsConditions:   termsConditions,
		ImageURL:          imageURL,
		IsExclusive:       isExclusive,
		Status:            "active",
		ViewsCount:        0,
		RedemptionsCount:  0,
	}

	err := s.businessRepo.CreateMerchantOffer(ctx, offer)
	if err != nil {
		s.logger.Errorf("Failed to create merchant offer: %v", err)
		return nil, fmt.Errorf("failed to create merchant offer: %w", err)
	}

	s.logger.Infof("Merchant offer created: %s", offer.ID)
	return offer, nil
}

func (s *BusinessService) GetMerchantOffer(ctx context.Context, id string) (*models.MerchantOffer, error) {
	return s.businessRepo.GetMerchantOfferByID(ctx, id)
}

func (s *BusinessService) ListMerchantOffers(ctx context.Context, limit, offset int) ([]*models.MerchantOffer, error) {
	return s.businessRepo.ListMerchantOffers(ctx, limit, offset)
}

func (s *BusinessService) UpdateMerchantOffer(ctx context.Context, id string, updates map[string]interface{}) error {
	offer, err := s.businessRepo.GetMerchantOfferByID(ctx, id)
	if err != nil {
		return err
	}
	if offer == nil {
		return fmt.Errorf("merchant offer not found")
	}

	// Apply updates
	if v, ok := updates["title"]; ok {
		offer.Title = v.(string)
	}
	if v, ok := updates["description"]; ok {
		offer.Description = v.(string)
	}
	if v, ok := updates["discount_percentage"]; ok {
		offer.DiscountPercentage = v.(*int)
	}
	if v, ok := updates["original_price"]; ok {
		offer.OriginalPrice = v.(*float64)
	}
	if v, ok := updates["offer_price"]; ok {
		offer.OfferPrice = v.(*float64)
	}
	if v, ok := updates["valid_from"]; ok {
		offer.ValidFrom = v.(time.Time)
	}
	if v, ok := updates["valid_until"]; ok {
		offer.ValidUntil = v.(time.Time)
	}
	if v, ok := updates["terms_conditions"]; ok {
		offer.TermsConditions = v.(*string)
	}
	if v, ok := updates["image_url"]; ok {
		offer.ImageURL = v.(*string)
	}
	if v, ok := updates["is_exclusive"]; ok {
		offer.IsExclusive = v.(bool)
	}
	if v, ok := updates["status"]; ok {
		offer.Status = v.(string)
	}

	return s.businessRepo.UpdateMerchantOffer(ctx, offer)
}

// Sponsorships
func (s *BusinessService) CreateSponsorship(ctx context.Context, userID, title, description, sponsorshipType string, targetAmount float64, startDate time.Time, endDate *time.Time, beneficiary, imageURL *string) (*models.Sponsorship, error) {
	sponsorship := &models.Sponsorship{
		ID:              uuid.New().String(),
		UserID:          userID,
		Title:           title,
		Description:     description,
		SponsorshipType: sponsorshipType,
		TargetAmount:    targetAmount,
		CurrentAmount:   0,
		StartDate:       startDate,
		EndDate:         endDate,
		Status:          "active",
		Beneficiary:     beneficiary,
		ImageURL:        imageURL,
	}

	err := s.businessRepo.CreateSponsorship(ctx, sponsorship)
	if err != nil {
		s.logger.Errorf("Failed to create sponsorship: %v", err)
		return nil, fmt.Errorf("failed to create sponsorship: %w", err)
	}

	s.logger.Infof("Sponsorship created: %s", sponsorship.ID)
	return sponsorship, nil
}

func (s *BusinessService) GetSponsorship(ctx context.Context, id string) (*models.Sponsorship, error) {
	return s.businessRepo.GetSponsorshipByID(ctx, id)
}

func (s *BusinessService) ListSponsorships(ctx context.Context, limit, offset int) ([]*models.Sponsorship, error) {
	return s.businessRepo.ListSponsorships(ctx, limit, offset)
}

func (s *BusinessService) ContributeToSponsorship(ctx context.Context, sponsorshipID, userID string, amount float64, message *string, isAnonymous bool) (*models.SponsorshipContribution, error) {
	contribution := &models.SponsorshipContribution{
		ID:            uuid.New().String(),
		SponsorshipID: sponsorshipID,
		UserID:        userID,
		Amount:        amount,
		Message:       message,
		IsAnonymous:   isAnonymous,
	}

	err := s.businessRepo.CreateSponsorshipContribution(ctx, contribution)
	if err != nil {
		s.logger.Errorf("Failed to create sponsorship contribution: %v", err)
		return nil, fmt.Errorf("failed to create sponsorship contribution: %w", err)
	}

	// Update sponsorship amount
	_ = s.businessRepo.UpdateSponsorshipAmount(ctx, sponsorshipID, amount)

	s.logger.Infof("Sponsorship contribution created: %s", contribution.ID)
	return contribution, nil
}

func (s *BusinessService) GetSponsorshipContributions(ctx context.Context, sponsorshipID string) ([]*models.SponsorshipContribution, error) {
	return s.businessRepo.ListSponsorshipContributions(ctx, sponsorshipID)
}

// Escrow Transactions
func (s *BusinessService) CreateEscrowTransaction(ctx context.Context, businessID, buyerID, sellerID string, amount float64, description string, releaseConditions *string) (*models.EscrowTransaction, error) {
	transaction := &models.EscrowTransaction{
		ID:                uuid.New().String(),
		BusinessID:        businessID,
		BuyerID:           buyerID,
		SellerID:          sellerID,
		Amount:            amount,
		Description:       description,
		Status:            "pending",
		ReleaseConditions: releaseConditions,
	}

	err := s.businessRepo.CreateEscrowTransaction(ctx, transaction)
	if err != nil {
		s.logger.Errorf("Failed to create escrow transaction: %v", err)
		return nil, fmt.Errorf("failed to create escrow transaction: %w", err)
	}

	s.logger.Infof("Escrow transaction created: %s", transaction.ID)
	return transaction, nil
}

func (s *BusinessService) GetEscrowTransaction(ctx context.Context, id string) (*models.EscrowTransaction, error) {
	return s.businessRepo.GetEscrowTransactionByID(ctx, id)
}

func (s *BusinessService) GetUserEscrowTransactions(ctx context.Context, userID string, role string) ([]*models.EscrowTransaction, error) {
	return s.businessRepo.ListEscrowTransactions(ctx, userID, role)
}

func (s *BusinessService) FundEscrowTransaction(ctx context.Context, id string) error {
	transaction, err := s.businessRepo.GetEscrowTransactionByID(ctx, id)
	if err != nil {
		return err
	}
	if transaction == nil {
		return fmt.Errorf("transaction not found")
	}

	now := time.Now()
	transaction.Status = "funded"
	transaction.FundedAt = &now

	return s.businessRepo.UpdateEscrowTransaction(ctx, transaction)
}

func (s *BusinessService) ReleaseEscrowTransaction(ctx context.Context, id string) error {
	transaction, err := s.businessRepo.GetEscrowTransactionByID(ctx, id)
	if err != nil {
		return err
	}
	if transaction == nil {
		return fmt.Errorf("transaction not found")
	}

	now := time.Now()
	transaction.Status = "completed"
	transaction.CompletedAt = &now

	return s.businessRepo.UpdateEscrowTransaction(ctx, transaction)
}

func (s *BusinessService) CancelEscrowTransaction(ctx context.Context, id string) error {
	transaction, err := s.businessRepo.GetEscrowTransactionByID(ctx, id)
	if err != nil {
		return err
	}
	if transaction == nil {
		return fmt.Errorf("transaction not found")
	}

	now := time.Now()
	transaction.Status = "cancelled"
	transaction.CancelledAt = &now

	return s.businessRepo.UpdateEscrowTransaction(ctx, transaction)
}
