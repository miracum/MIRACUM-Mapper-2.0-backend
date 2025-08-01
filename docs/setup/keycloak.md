# Keycloak Setup

Regardless of whether an existing Keycloak server or a new server is used, a Keycloak client for the miracum mapper must be created.

Please login to the Keycloak main site with an admin account.
For Development and Testing, Keycloak should be available at `http://localhost:8081` and a temporary admin account exists with the username `admin` and password `admin`.

Import the miracum-mapper client by clicking on `Client` and then `Import Client`. Select the file which can be found [here](/tools/setup/keycloak-client-miracum-mapper.json). Afterward, roles need to be created in the section `Roles` within the newly created `miracum-mapper` client. Two roles should be created, called `admin` and `normal`. These roles need to be assigned to either the admin user or a newly created user in oder for them to access the miracum mapper. A new user can be created by clicking on `Users`, `Add User` and following the dialog. Afterwards, click on the user, then on `Role Mapping`, `Assign Role`, in the search bar search for either `admin` or `normal`, select the role and click on `assign`.
