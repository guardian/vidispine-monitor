# vidispine-monitor

## What is it?
Vidispine monitor is a small piece of software that performs systems
checks on an active Vidispine server at a set interval.

If problems are found, then alerts are sent to the PagerDuty monitoring
system for action by administrators.

## Checks carried out

Checks are carried out in order, not in parallel.  If an internal error is found
trying to carry out the check (e.g. Vidispine responded with a 500, or content could
not be parsed, Pagerduty is offline, or we ran out of memory etc.) then this message is displayed and the
app exits.  In Kubernetes this causes a crashloop state, which should be easy to
identify in order to isolate the issue.

Error exit only occurs after _every_ check has been completed.

### 1. System health
The /healthcheck/ endpoint on the 9001 admin port is checked for all subcomponents;
this includes message broker, index, database, etc.  A message is sent for every failure
identified by the server

### 2. Storages
The /API/storage endpoint on the 8080 API port is checked.  For each storage identified,
the state is checked as well as whether the storage is over the high watermark and
whether it is over 95% full.

This check requires regular API permissions to work, refer to README.md in the
`vidispine` subdirectory to see how to set these up.

## Build and deployment

You need to have Go installed, ideally version 1.14 or later (modules support
is a must).  You'll make life easier if you have GNU make installed too.

Assuming you have these two, simply:

```bash
$ make clean && make
$ docker build . -t your_org/vidispine-monitor:DEV
```

For a deployment manifest, refer to `vidispine/vidispine-monitor.yaml` in
the prexit-local repo.  You should only ever need one instance of this app
running at a time.