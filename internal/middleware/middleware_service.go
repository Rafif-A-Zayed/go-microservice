package middleware

import user "user-management/internal"

// Middleware describes a service middleware.
type Middleware func(service user.Service) user.Service
