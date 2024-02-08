# Product Catalog Demo App

## Cases Validation

In order to validate the initial issue and then validate the fix,
we have created a script that will run the cases and validate the results.

The script will execute the `fury get-token` command internally and then
will execute a k6 script to perform a small load test.

This script is located in the `validation.sh` file. It will run the cases
following the format:

```shell
./validation.sh <app_name> <group_number:1|2|3|4|5|6|7|8> <env:local|fury> <case_number:1|2|3|4>
```

For instance, considering the following parameters:

- app_name: `drs-devbego`
- group_number: `1`
- env: `fury`
- case_number: `1`

The execution command will be:

```shell
./validation.sh drs-devbego 1 fury 1
```

## Pre Requisites
- K6: `brew install k6`
