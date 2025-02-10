1) main.go
    1. Services (The Brain Surgeons)
        What They Do:
        Contain core business logic
        Handle data manipulation
        Interact with databases
        Perform complex operations
        Are reusable across different parts of the application

        When to Use Services:
        When working with databases
        When implementing business rules ("Users must have unique emails")
        When performing data transformations
        When handling transactions


    2. Handlers (The Receptionists)
        What They Do:
        Handle HTTP requests/responses
        Parse input data
        Call appropriate services
        Return proper HTTP status codes
        Don't contain business logic    

2) structure of a function
        func (s *AuthService) CreateUser(user *models.User) error {
            // Validate input
            if user.Username == "" || user.Password == "" || user.Email == "" {
                return errors.New("username, email, and password are required")
            }

            // Hash the password
            if err := user.HashPassword(); err != nil {
                return err
            }

            // Create user in database
            result := s.DB.Create(user)
            return result.Error
        }
        
        Receiver (s *AuthService):
        The "factory" where this work happens
        *AuthService pointer gives access to the factory's tools (database connection)
        Using a pointer receiver (*AuthService) instead of a value receiver (AuthService) means Go won't create a copy of the entire AuthService struct when the method is called
        If we used (s AuthService) instead:
            Every method call would create a complete copy of AuthService
            Any modifications to s would only affect the copy, not the original AuthService
            Memory inefficient if AuthService contains large fields like database connections
            Changes made within the method wouldn't persist

        Parameters (user *models.User):
        Raw materials arriving at the assembly line
        *models.User pointer avoids copying the whole user struct (efficient)
        The unfinished product moving through the line
        depth:
            callers memory would already created a User(), i.e check where this CreateUser function is called it would already have User() initialized some where at top.
            Caller's memory:             Function's memory:
                User{                    pointer (8 bytes) ----→ [Points to original User]
                ID: "123"              
                Username: "john"       
                ...
                }
            we are just reusing caller memory; we cn have multiple function work on same caller memory and at end do something with final value of user(eg. save it somewhere)    