# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.

SELF_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
CHECK_PLANFILE_PATH ?= check-plan.output

include $(SELF_DIR)/common.mk

all:
.PHONY: all

setup: ## set up local dependencies for this repo
	$(MAKE) -C $(REPO_ROOT) setup
.PHONY: setup

check: lint check-plan ## run all checks for this component
.PHONY: check

fmt: terraform ## format code in this component
	$(terraform_command) fmt $(TF_ARGS)
.PHONY: fmt

lint: lint-terraform-fmt lint-tflint ## run all linters for this component
.PHONY: lint

lint-tflint: ## run the tflint linter for this component
	@printf "tflint: "
ifeq ($(TFLINT_ENABLED),1)
	@tflint || exit $$?;
else
	@echo "disabled"
endif
.PHONY: lint-tflint

lint-terraform-fmt: terraform ## run `terraform fmt` in check mode
	@printf "fmt check: "
	@for f in $(TF); do \
		printf . \
		$(terraform_command) fmt $(TF_ARGS) --check=true --diff=true $$f || exit $$? ; \
	done
	@echo
.PHONY: lint-terraform-fmt

check-auth: ## check that authentication is properly set up for this component
	@for p in $(AWS_BACKEND_PROFILE) $(AWS_PROVIDER_PROFILE); do \
		aws --profile $$p sts get-caller-identity > /dev/null || echo "AWS AUTH error. This component is configured to use a profile name $$p. Please add one to your ~/.aws/config"; \
	done
.PHONY: check-auth

ifeq ($(MODE),local)
plan: check-auth init fmt ## run a terraform plan
	@$(terraform_command) plan $(TF_ARGS) -refresh=$(REFRESH_ENABLED) -input=false
else ifeq ($(MODE),atlantis)
plan: check-auth init lint
	@$(terraform_command) plan $(TF_ARGS) -refresh=$(REFRESH_ENABLED) -input=false -out $(PLANFILE) | scenery
else
	@echo "Unknown MODE: $(MODE)"
	@exit -1
endif
.PHONY: plan

ifeq ($(MODE),local)
apply: check-auth init ## run a terraform apply
ifeq ($(ATLANTIS_ENABLED),1)
ifneq ($(FORCE),1)
	@echo "`tput bold`This component is configured to be used with atlantis. If you want to override and apply locally, add FORCE=1.`tput sgr0`"
	exit -1
endif
endif
	@$(terraform_command) apply $(TF_ARGS) -refresh=$(REFRESH_ENABLED) -auto-approve=$(AUTO_APPROVE)
else ifeq ($(MODE),atlantis)
apply: check-auth ## run a terraform apply
	@$(terraform_command) apply $(TF_ARGS) -refresh=$(REFRESH_ENABLED) -auto-approve=true $(PLANFILE)
else
	echo "Unknown mode: $(MODE)"
	exit -1
endif
.PHONY: apply

docs:
	echo
.PHONY: docs

clean: ## clean modules and plugins for this component
	-rm -rfv .terraform/modules
	-rm -rfv .terraform/plugins
.PHONY: clean

test:
.PHONY: test

init: terraform check-auth ## run terraform init for this component
ifeq ($(MODE),local)
	@$(terraform_command) init -input=false
else ifeq ($(MODE),atlantis)
	@$(terraform_command) init $(TF_ARGS) -input=false
else
	@echo "Unknown MODE: $(MODE)"
	@exit -1
endif
.PHONY: init

check-plan: init check-auth ## run a terraform plan and check that it does not fail
	@$(terraform_command) plan $(TF_ARGS) -detailed-exitcode -lock=false -out=$(CHECK_PLANFILE_PATH) ; \
	ERR=$$?; \
	if [ $$ERR -eq 0 ] ; then \
		echo "Success"; \
	elif [ $$ERR -eq 1 ] ; then \
		echo "Error in plan execution."; \
		exit 1; \
	elif [ $$ERR -eq 2 ] ; then \
		echo "Diff";  \
	fi ; \
	if [ -n "$(BUILDEVENT_FILE)" ]; then \
		fogg exp entropy -f $(CHECK_PLANFILE_PATH) -o $(BUILDEVENT_FILE) ; \
	fi
	rm $(CHECK_PLANFILE_PATH)
.PHONY: check-plan

run: check-auth ## run an arbitrary terraform command, CMD. ex `make run CMD='show'`
	@$(terraform_command) $(CMD)
.PHONY: run
