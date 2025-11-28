# Email Verification Implementation Plan

## Overview
Implement token-based email verification for user accounts using SendGrid for email delivery.

---

## Phase 1: Database Schema Updates

### Goal
Add email verification fields to the User model and migrate the database.

### Tasks
1. **Update User model** (`server/main.go:92`)
   - Add `EmailVerified bool` field (default: false)
   - Add `VerificationToken *string` field (nullable, indexed)
   - Add `TokenExpiresAt *time.Time` field (nullable)

2. **Test database migration**
   - Run server to trigger auto-migration
   - Verify new columns exist in database
   - Ensure existing users are handled (EmailVerified defaults to false)

### Acceptance Criteria
- [ ] User model has three new fields
- [ ] Database migration completes successfully
- [ ] Existing users can still log in

---

## Phase 2: Token Generation & Validation

### Goal
Create secure utilities for generating and validating verification tokens.

### Tasks
1. **Add token generation function**
   - Use `crypto/rand` for secure random token generation
   - Encode as URL-safe base64 string (32 bytes → ~44 characters)
   - Add function near other utilities (around line 43)

2. **Add token validation helper**
   - Check token exists and matches database
   - Verify token hasn't expired
   - Return user if valid, error otherwise

3. **Set token expiration policy**
   - Default: 24 hours from generation
   - Make configurable via constant for easy adjustment

### Acceptance Criteria
- [ ] Tokens are cryptographically random
- [ ] Tokens are URL-safe
- [ ] Validation checks expiration correctly
- [ ] Unit tests pass for token generation/validation

---

## Phase 3: SendGrid Integration

### Goal
Set up email delivery service and create email templates.

### Tasks
1. **Install SendGrid Go SDK**
   - Run: `cd server && go get github.com/sendgrid/sendgrid-go`
   - Add to imports in main.go

2. **Configure SendGrid API key**
   - Add `SENDGRID_API_KEY` to `server/.env`
   - Get API key from SendGrid dashboard (free tier: 100 emails/day)

3. **Create email sending function**
   - Function: `SendVerificationEmail(toEmail, token string) error`
   - Build verification URL: `http://localhost:5173/verify?token={token}`
   - Create HTML email template with:
     - Welcome message
     - Clear "Verify Email" button/link
     - Token expiration warning (24 hours)
     - Fallback plain text version

4. **Add email configuration**
   - Sender email and name
   - Support for environment-based URLs (dev vs production)

### Acceptance Criteria
- [ ] SendGrid SDK installed and imported
- [ ] API key loaded from environment
- [ ] Email function sends successfully
- [ ] Email template is professional and clear
- [ ] Links use correct base URL for environment

---

## Phase 4: Update Registration Flow

### Goal
Modify the `/register` endpoint to generate tokens and send verification emails.

### Tasks
1. **Update `/register` endpoint** (`server/main.go:333`)
   - After creating user successfully:
     - Generate verification token
     - Set token expiration (24 hours from now)
     - Save token to user record in database
     - Send verification email
   - Handle email sending errors gracefully:
     - Log error but don't block registration
     - User can request new verification email later

2. **Update registration response**
   - Include message: "Please check your email to verify your account"
   - Still return user data and potentially JWT (see Phase 6 for policy decision)

3. **Add error handling**
   - If email fails to send, log error but complete registration
   - Return warning in response if email failed

### Acceptance Criteria
- [ ] New users receive verification email immediately
- [ ] Token is saved to database correctly
- [ ] Registration succeeds even if email fails
- [ ] User receives clear instructions in response

---

## Phase 5: Email Verification Endpoint

### Goal
Create endpoint for users to verify their email via token.

### Tasks
1. **Create `/verify-email` POST endpoint** (public route)
   - Accept JSON body: `{"token": "abc123..."}`
   - Validate token format
   - Look up user by token
   - Check token expiration
   - If valid:
     - Set `EmailVerified = true`
     - Clear `VerificationToken` and `TokenExpiresAt`
     - Return success message
   - If invalid/expired:
     - Return appropriate error message

2. **Create `/resend-verification` POST endpoint** (requires auth)
   - Check if user is already verified
   - Generate new token
   - Send new verification email
   - Rate limit: prevent spam (max 1 email per 5 minutes per user)

### Acceptance Criteria
- [ ] Valid tokens successfully verify accounts
- [ ] Expired tokens return clear error message
- [ ] Invalid tokens return security-safe error (don't leak info)
- [ ] Users can request new verification emails
- [ ] Rate limiting prevents email spam

---

## Phase 6: Access Control & Middleware

### Goal
Restrict unverified users from accessing certain features.

### Tasks
1. **Create `RequireVerifiedEmail()` middleware**
   - Check if current user's `EmailVerified` is true
   - If false, return 403 Forbidden with message
   - Place after `RequireRole()` in middleware chain

2. **Decide verification policy**
   - **Option A (Strict)**: Require verification for all authenticated routes
   - **Option B (Lenient)**: Only require for critical actions (surveys, responses)
   - **Recommended**: Option B for better UX

3. **Apply middleware to routes**
   - Student routes:
     - `/student/responses` (POST) - require verification
     - `/student/surveys` (GET) - allow unverified (can browse)
   - Professor routes:
     - `/professor/surveys` (POST) - require verification
   - Admin routes:
     - Consider exempting admins or manual verification

4. **Update `/login` response**
   - Include `email_verified` status in user object
   - Frontend can show verification banner if false

### Acceptance Criteria
- [ ] Unverified users receive clear error when accessing protected routes
- [ ] Verified users have normal access
- [ ] Login response includes verification status
- [ ] Policy is consistently applied across routes

---

## Phase 7: Frontend Integration

### Goal
Create UI flows for email verification in SvelteKit app.

### Tasks
1. **Create verification page** (`client/src/routes/verify/+page.svelte`)
   - Extract token from URL query parameter
   - Auto-submit verification request on page load
   - Show loading state → success/error message
   - Redirect to dashboard on success

2. **Add verification banner to dashboard**
   - Show if `user.email_verified === false`
   - Include "Resend Verification Email" button
   - Display countdown/rate limit if recently sent

3. **Update registration page**
   - Show success message: "Check your email to verify your account"
   - Provide link to resend verification email

4. **Handle verification errors in UI**
   - Expired token → show resend button
   - Invalid token → show error, redirect to login
   - Network errors → retry logic

### Acceptance Criteria
- [ ] Users can verify email by clicking link
- [ ] Verification page handles all states (loading, success, error)
- [ ] Unverified users see clear prompts to verify
- [ ] Resend functionality works from UI
- [ ] UX is smooth and non-blocking

---

## Phase 8: Testing & Validation

### Goal
Ensure email verification system is robust and secure.

### Tasks
1. **Write backend tests** (`server/main_test.go`)
   - Token generation produces unique values
   - Token validation checks expiration correctly
   - Verification endpoint validates tokens properly
   - Middleware blocks unverified users
   - Resend endpoint enforces rate limiting

2. **Manual testing scenarios**
   - Register new user → receive email → verify → access protected routes
   - Try expired token → see error → resend → verify successfully
   - Try invalid token → see error message
   - Verified user can access all features
   - Unverified user blocked from protected routes

3. **Security checks**
   - Tokens are cryptographically random
   - Tokens expire after 24 hours
   - Used tokens are cleared from database
   - Rate limiting prevents abuse
   - Error messages don't leak user information

### Acceptance Criteria
- [ ] All unit tests pass
- [ ] Manual test scenarios complete successfully
- [ ] No security vulnerabilities identified
- [ ] Edge cases handled gracefully

---

## Phase 9: Production Preparation

### Goal
Prepare email verification for production deployment.

### Tasks
1. **Environment configuration**
   - Document required environment variables:
     - `SENDGRID_API_KEY`
     - `FRONTEND_URL` (for verification links)
     - `SENDER_EMAIL`, `SENDER_NAME`
   - Add production SendGrid API key (move from free to paid tier if needed)

2. **Email template improvements**
   - Use production domain in verification links
   - Add branding (logo, colors)
   - Ensure mobile-responsive design
   - Include support contact info

3. **Migration strategy for existing users**
   - Existing users will have `EmailVerified = false` after migration
   - Options:
     - **Auto-verify**: Set all existing users to verified
     - **Require verification**: Send bulk verification emails
     - **Grandfather**: Allow unverified access for existing users
   - **Recommended**: Auto-verify existing users

4. **Monitoring & logging**
   - Log email sending failures
   - Track verification rates (how many users verify?)
   - Alert on high failure rates

5. **Documentation**
   - Update README with email verification setup instructions
   - Document environment variables
   - Add troubleshooting guide for common email issues

### Acceptance Criteria
- [ ] Production environment configured
- [ ] Email templates are production-ready
- [ ] Existing user migration strategy decided and documented
- [ ] Monitoring in place
- [ ] Documentation complete

---

## Success Metrics

- **Security**: All new accounts require email verification
- **User Experience**: Verification process is smooth and clear
- **Reliability**: Emails deliver consistently (>95% success rate)
- **Performance**: No significant impact on registration time
- **Maintainability**: Code is well-tested and documented

---

## Rollback Plan

If issues arise after deployment:
1. Disable verification middleware (remove `RequireVerifiedEmail()` calls)
2. Continue collecting verification data without enforcing
3. Fix issues and re-enable gradually (start with new users only)
