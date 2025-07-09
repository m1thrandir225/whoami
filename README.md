# whoami (WIP)

**This project is actively under development and is a work-in-progress**.

A Central Authentication Service written entirely in Go. The service is intented to be fully Open-Source and used as a seperate service in your application. I plan to make it easily plug-and-play for microservices architectures but also for maybe just frontend applications that to use some authentication features instead of using something like Firebase or Supabase.

## Features

- Basic Username/Password Authentication
- Email Verification
- OAuth2 via supported providers
- Password Reset Flows
- Account activation/deactivation
- Rate limiting per ip/user
- Suspicious activity detection
- Password History
- Short lived access-tokens and longer refresh tokens
- Password Strength requirements
- HaveIBeenPwned integration
- Account lockout after multiple failed login attempts

## Architecture (WIP)

![Database Design](./.github/images/whoami-db.png)

## OAuth 2 Supported Providers

The following are the planned supported OAuth 2 providers:

- Github
