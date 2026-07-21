-- Business Listings
CREATE TABLE IF NOT EXISTS business_listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    business_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    location VARCHAR(255) NOT NULL,
    website VARCHAR(500),
    description TEXT NOT NULL,
    instagram_handle VARCHAR(100),
    facebook_handle VARCHAR(100),
    logo_url VARCHAR(500),
    banner_url VARCHAR(500),
    is_verified BOOLEAN DEFAULT FALSE,
    is_featured BOOLEAN DEFAULT FALSE,
    featured_until TIMESTAMP,
    status VARCHAR(50) DEFAULT 'active', -- active, suspended, deleted
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_business_listings_user_id ON business_listings(user_id);
CREATE INDEX idx_business_listings_status ON business_listings(status);
CREATE INDEX idx_business_listings_featured ON business_listings(is_featured, featured_until);
CREATE INDEX idx_business_listings_location ON business_listings(location);

-- Jobs
CREATE TABLE IF NOT EXISTS jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    business_id UUID REFERENCES business_listings(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    requirements TEXT,
    responsibilities TEXT,
    location VARCHAR(255),
    job_type VARCHAR(50) DEFAULT 'full-time', -- full-time, part-time, contract, remote, internship
    salary_range VARCHAR(100),
    application_deadline TIMESTAMP,
    status VARCHAR(50) DEFAULT 'active', -- active, closed, draft
    views_count INTEGER DEFAULT 0,
    applications_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_jobs_user_id ON jobs(user_id);
CREATE INDEX idx_jobs_business_id ON jobs(business_id);
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_type ON jobs(job_type);
CREATE INDEX idx_jobs_deadline ON jobs(application_deadline);

-- Job Applications
CREATE TABLE IF NOT EXISTS job_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cover_letter TEXT,
    resume_url VARCHAR(500),
    status VARCHAR(50) DEFAULT 'pending', -- pending, reviewed, shortlisted, rejected, hired
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(job_id, user_id)
);

CREATE INDEX idx_job_applications_job_id ON job_applications(job_id);
CREATE INDEX idx_job_applications_user_id ON job_applications(user_id);
CREATE INDEX idx_job_applications_status ON job_applications(status);

-- Tenders
CREATE TABLE IF NOT EXISTS tenders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    requirements TEXT,
    budget_range VARCHAR(100),
    submission_deadline TIMESTAMP NOT NULL,
    tender_number VARCHAR(100),
    category VARCHAR(100),
    status VARCHAR(50) DEFAULT 'open', -- open, closed, awarded, cancelled
    views_count INTEGER DEFAULT 0,
    bids_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tenders_user_id ON tenders(user_id);
CREATE INDEX idx_tenders_status ON tenders(status);
CREATE INDEX idx_tenders_deadline ON tenders(submission_deadline);
CREATE INDEX idx_tenders_category ON tenders(category);

-- Tender Bids
CREATE TABLE IF NOT EXISTS tender_bids (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tender_id UUID NOT NULL REFERENCES tenders(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_name VARCHAR(255) NOT NULL,
    proposal_description TEXT NOT NULL,
    bid_amount DECIMAL(15, 2),
    proposal_url VARCHAR(500),
    status VARCHAR(50) DEFAULT 'submitted', -- submitted, under_review, awarded, rejected
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tender_id, user_id)
);

CREATE INDEX idx_tender_bids_tender_id ON tender_bids(tender_id);
CREATE INDEX idx_tender_bids_user_id ON tender_bids(user_id);
CREATE INDEX idx_tender_bids_status ON tender_bids(status);

-- Class Groups (SBC & SGC)
CREATE TABLE IF NOT EXISTS class_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_type VARCHAR(50) NOT NULL, -- SBC, SGC
    year_of_completion INTEGER NOT NULL,
    class_name VARCHAR(100) NOT NULL, -- e.g., "Form 1 2020"
    description TEXT,
    class_rep_id UUID REFERENCES users(id) ON DELETE SET NULL,
    member_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(school_type, year_of_completion, class_name)
);

CREATE INDEX idx_class_groups_school_type ON class_groups(school_type);
CREATE INDEX idx_class_groups_year ON class_groups(year_of_completion);
CREATE INDEX idx_class_groups_active ON class_groups(is_active);

-- Class Group Members
CREATE TABLE IF NOT EXISTS class_group_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_group_id UUID NOT NULL REFERENCES class_groups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'member', -- member, admin, class_rep
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(class_group_id, user_id)
);

CREATE INDEX idx_class_group_members_group_id ON class_group_members(class_group_id);
CREATE INDEX idx_class_group_members_user_id ON class_group_members(user_id);
CREATE INDEX idx_class_group_members_role ON class_group_members(role);

-- Merchant Offers
CREATE TABLE IF NOT EXISTS merchant_offers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    business_id UUID REFERENCES business_listings(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    discount_percentage INTEGER,
    original_price DECIMAL(10, 2),
    offer_price DECIMAL(10, 2),
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    terms_conditions TEXT,
    image_url VARCHAR(500),
    is_exclusive BOOLEAN DEFAULT FALSE,
    status VARCHAR(50) DEFAULT 'active', -- active, expired, suspended
    views_count INTEGER DEFAULT 0,
    redemptions_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_merchant_offers_user_id ON merchant_offers(user_id);
CREATE INDEX idx_merchant_offers_business_id ON merchant_offers(business_id);
CREATE INDEX idx_merchant_offers_status ON merchant_offers(status);
CREATE INDEX idx_merchant_offers_validity ON merchant_offers(valid_from, valid_until);

-- Sponsorships/Endowments
CREATE TABLE IF NOT EXISTS sponsorships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    sponsorship_type VARCHAR(50) NOT NULL, -- education, sports, infrastructure, events, other
    target_amount DECIMAL(15, 2) NOT NULL,
    current_amount DECIMAL(15, 2) DEFAULT 0,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    status VARCHAR(50) DEFAULT 'active', -- active, completed, cancelled, paused
    beneficiary VARCHAR(255),
    image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sponsorships_user_id ON sponsorships(user_id);
CREATE INDEX idx_sponsorships_type ON sponsorships(sponsorship_type);
CREATE INDEX idx_sponsorships_status ON sponsorships(status);
CREATE INDEX idx_sponsorships_dates ON sponsorships(start_date, end_date);

-- Sponsorship Contributions
CREATE TABLE IF NOT EXISTS sponsorship_contributions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sponsorship_id UUID NOT NULL REFERENCES sponsorships(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(15, 2) NOT NULL,
    message TEXT,
    is_anonymous BOOLEAN DEFAULT FALSE,
    contributed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sponsorship_contributions_sponsorship_id ON sponsorship_contributions(sponsorship_id);
CREATE INDEX idx_sponsorship_contributions_user_id ON sponsorship_contributions(user_id);

-- Escrow Transactions (for business listings)
CREATE TABLE IF NOT EXISTS escrow_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL REFERENCES business_listings(id) ON DELETE CASCADE,
    buyer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seller_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(15, 2) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'pending', -- pending, funded, in_progress, completed, released, cancelled, disputed
    release_conditions TEXT,
    funded_at TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_escrow_transactions_business_id ON escrow_transactions(business_id);
CREATE INDEX idx_escrow_transactions_buyer_id ON escrow_transactions(buyer_id);
CREATE INDEX idx_escrow_transactions_seller_id ON escrow_transactions(seller_id);
CREATE INDEX idx_escrow_transactions_status ON escrow_transactions(status);

-- Update triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_business_listings_updated_at BEFORE UPDATE ON business_listings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_jobs_updated_at BEFORE UPDATE ON jobs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_job_applications_updated_at BEFORE UPDATE ON job_applications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tenders_updated_at BEFORE UPDATE ON tenders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tender_bids_updated_at BEFORE UPDATE ON tender_bids
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_class_groups_updated_at BEFORE UPDATE ON class_groups
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_merchant_offers_updated_at BEFORE UPDATE ON merchant_offers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sponsorships_updated_at BEFORE UPDATE ON sponsorships
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_escrow_transactions_updated_at BEFORE UPDATE ON escrow_transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
