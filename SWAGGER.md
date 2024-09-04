**Note:** client will generate schema from swagger docs so, please follow these configurations

- ### Struct

    - use **validate:"required"** tag to insure generated fields/keys are not null/optional.
      ```go
      type User struct {
          Name string  `json:"name" validate:"required"`
      }
      ```
      will be generated as
      ```ts
      interface User {
        name: string;
      }
      ```
      else
      ```ts
      interface User {
        name?: string;
      }
      ```

- ### Function Comments

    - Using `// @Tags` for User and Admin Access endpoints
        - For easy navigation and understanding, utilize distinct @Tags values for endpoints accessible by users and admins. For instance,
        - use:
            - `// @Tags UserApi` for endpoints under /users accessible to regular users.
            - `// @Tags UserManagementApi` for endpoints under /users restricted to administrative access.
