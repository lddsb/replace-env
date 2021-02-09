### replace-env
It works great when there has different `.env`(`.env.json`) file for different environment.

### How it works
You need set your different environment variables for different environment in your operate system.For example:
```shell
...
export DEV_HTTPS=false
export PROD_HTTPS=true
...
```

and your `.env` should like:
```
...
HTTPS=
...
```

### Usage
```shell
replace-env source-file output-file
```

#### .env
```shell
replace-env e .env.example .env
```

#### .env.json
```shell
replace j .env.json.example .env.json
```

#### custom branch environment variables(default: `CI_COMMIT_BRANCH`)
##### Drone CI
```shell
replace e --branch-env DRONE_BRANCH .env.example .env
```

#### GitLab CI
```shell
replace e --branch-env CI_COMMIT_BRANCH .env.example .env
```
