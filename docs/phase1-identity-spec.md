# Phase 1 — Identity / Collaborators

Evidence baseline: `docs/reverse-spec/account.md` (`CONFIRMED_AUTH`).

## Domain decisions
- `Cargo` is descriptive metadata, not authorization.
- Authorization uses an explicit permission set per tenant membership.
- A user identity can belong to multiple tenants through memberships.
- Membership status and session revocation are independent operations.
- Email is normalized to lowercase.
- CPF is normalized to 11 digits and validated.
- Passwords must satisfy the visible target policy before hashing.

## Initial acceptance criteria
1. Collaborator requires name, CPF, email, phone and title.
2. Membership requires at least one valid permission.
3. Password policy requires >= 8 chars, upper, lower, digit and symbol.
4. Unknown permission codes are rejected.
5. Duplicate permission codes are canonicalized.
6. New memberships start active.
7. Deactivation does not imply session revocation.
