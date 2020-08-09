Backend
=======

# Code Structure

1. `cmd/` contains program entry point.

1. `database/` contains database migrations and the relevant scripts.

1. `docker/` docker related configurations.

1. `internal/` code that is called by program entrypoint. If it's HTTP Api it contains server and routes definitions.

    Code in this folder should not make sense to be exported and be used by other packages.

1. `pkg/` code which are exported and can be used by other packages.

    `pkg/core` contains our business logic.

    Other packages are library code each handles specific thing.
