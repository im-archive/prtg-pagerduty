# PRTG/PagerDuty Notification Integration


[![Build Status](https://travis-ci.org/WTW-IM/prtg-pagerduty.svg?branch=master)](https://travis-ci.org/WTW-IM/prtg-pagerduty)
[![Go Report Card](https://goreportcard.com/badge/github.com/WTW-IM/prtg-pagerduty)](https://goreportcard.com/report/github.com/WTW-IM/prtg-pagerduty)

## Goals

* Create incidents using version 2 of the PagerDuty Events API for triggered PRTG alerts.

* Automatically resolve alerts when status returns to normal or paused in PRTG.


## Build & Installation

Build the package

`go get github.com/WTW-IM/prg-pagerduty`

`go build`

From an Adminstrator powershell session:

`cp pagerduty.exe "C:\Program Files (x86)\PRTG Network Monitor\Notifications\EXE\"`


## Configuring notification in PRTG

Create new basic notification. Check "EXECUTE PROGRAM" selecting pagerduty.exe from the Program File dropdown.

Populate the parameter field with the following, substituting the service key with your service integration key

`-probe "%probe" -device "%device" -name "%name" -status "%status" -date "%datetime" -linkdevice %linkdevice -message "%message" -servicekey myShineyV2IntegrationKey -severity "critical"`
