# Features

# Roles:

## Application Roles

There are roles which can be set in KeyCloak to set permissions concerning the whole Miracum Mapper application. The roles are `admin` and `user`.

- `user`: The user role is the default role for all users. It allows to access the application and see all projects where the user is a member of. What the user can do in the project is defined by the project roles which are described later on
- `admin`: The admin role is a special role which has all permissions of the user and can additionally create new projects and assign users to projects. Respectively the admin can also adjust or remove users from projects or edit or delete projects. The creation, editing and deletion of Code Systems is also only allowed for users with the admin role.

When the roles should have other names in the keycloak instance, there are two places where the name has to be adjusted:

1. In [internal/config/constants.go](../internal/config/constants.go) `KeycloakAdminScope` and `KeycloakUserScope` have to be adjusted to the new role names
2. In the [api/openapi.yaml](../api/openapi.yaml) the `securitySchemes` have to be adjusted to the new role names

Then a new docker container has to be built in order for the changes to take effect.

## Project Roles

There are roles which can be set for each project individually to restrict the actions a user can do in the project. The roles are `reviewer`, `editor` and `project_owner`.

- `reviewer`: The reviewer role is the role with least privileges. The reviewer can only view the project and the mappings within them. The creation of new mapping and the deletion of existing mappings is not allowed for the reviewer. Only comments, status and equivalence of a mapping can be adjusted by the reviewer. Therefore the concepts can't be changed. This role should get extended further in the future to allow more actions for the reviewer when proper version management of the mappings is implemented. A look at MIRACUM Mapper 1 can be helpful to see what actions the reviewer should be able to do.
- `editor`: The editor role can do everything the reviewer can do and additionally create new mappings and delete existing mappings. He is also able to fully edit mappings. The editor can not delete the project or adjust the project settings.
- `project_owner`: The project owner is the role with the most privileges. The project owner can do everything the editor can do and additionally adjust the project settings. The project owner can also assign roles to other users in the project.

## Additional Information

The rules are just an initial implementation for the mapper and should be extended further in the future. When a proper workflow management is implemented, new capabilities or new roles in general can be added.
