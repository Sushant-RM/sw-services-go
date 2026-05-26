#!/bin/bash
# full_demo_runner.sh: Coordinates and monitors the entire end-to-end municipal validation suite

# Grant execution permissions
chmod +x scratch/*.sh

clear
echo -e "\033[1;36m======================================================================\033[0m"
echo -e "\033[1;36m             DIGIT MUNICIPAL WORKFLOW MIGRATION SUITE                 \033[0m"
echo -e "\033[1;36m======================================================================\033[0m"
echo -e "Started     : $(date '+%Y-%m-%d %H:%M:%S')"
echo -e "Target Node : sw-services-go microservice"
echo -e "Network     : digit_sw_net (Docker-Compose orchestrated environment)"
echo -e "======================================================================"

run_suite() {
  local script=$1
  local name=$2
  echo -e "\nRunning: $name..."
  bash "$script"
  local status=$?
  if [ $status -eq 0 ]; then
    echo -e "\033[1;32mSUCCESS: $name [PASS]\033[0m"
    return 0
  else
    echo -e "\033[1;31mFAILURE: $name [FAIL]\033[0m"
    return 1
  fi
}

ERRORS=0

run_suite "scratch/ecosystem_check.sh" "Ecosystem Check" || ERRORS=$((ERRORS+1))
run_suite "scratch/citizen_flow.sh" "Citizen Ingestion & Submission Flow" || ERRORS=$((ERRORS+1))
run_suite "scratch/verifier_flow.sh" "Verifier & Inspector Back-Office Flow" || ERRORS=$((ERRORS+1))
run_suite "scratch/approval_flow.sh" "Approval, Payment & Activation Flow" || ERRORS=$((ERRORS+1))
run_suite "scratch/db_audit.sh" "Database Direct DDL/DML Verification Audit" || ERRORS=$((ERRORS+1))
run_suite "scratch/kafka_audit.sh" "Kafka Stream Event Audit Verification" || ERRORS=$((ERRORS+1))
run_suite "scratch/edgecase_validation.sh" "API Exceptions & Failure Resiliency Probes" || ERRORS=$((ERRORS+1))

echo -e "\n======================================================================"
echo -e "                      MIGRATION SUITE SUMMARY"
echo -e "======================================================================"
if [ $ERRORS -eq 0 ]; then
  echo -e "  * Total Sub-Suites Run : 7"
  echo -e "  * Total Successes      : 7"
  echo -e "  * Total Failures       : 0"
  echo -e "\n  STATUS VERIFICATION: \033[1;32m[ALL PASSED]\033[0m Go-converted microservice ready for review."
  exit 0
else
  echo -e "  * Total Sub-Suites Run : 7"
  echo -e "  * Total Successes      : $((7-ERRORS))"
  echo -e "  * Total Failures       : $ERRORS"
  echo -e "\n  STATUS VERIFICATION: \033[1;31m[FAILED SUITE]\033[0m Please inspect container logs for failed transactions."
  exit 1
fi
