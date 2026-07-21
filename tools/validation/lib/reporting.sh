#!/usr/bin/env bash

validation_report_init() {
    local title="$1"
    local report_dir="$2"

    VALIDATION_REPORT_TITLE="$title"
    VALIDATION_REPORT_DIR="$report_dir"
    mkdir -p "$VALIDATION_REPORT_DIR"

    VALIDATION_CHECK_NAMES=()
    VALIDATION_CHECK_STATUS=()
    VALIDATION_CHECK_EXIT=()
    VALIDATION_CHECK_COMMAND=()
    VALIDATION_CHECK_LOG=()
    VALIDATION_CHECK_CAUSE=()
    VALIDATION_SKIP_NAMES=()
    VALIDATION_SKIP_REASONS=()
    VALIDATION_REPORT_NOTES=()

    VALIDATION_PASS_COUNT=0
    VALIDATION_FAIL_COUNT=0
    VALIDATION_SKIP_COUNT=0
}

validation_note() {
    local note="$*"
    VALIDATION_REPORT_NOTES+=("$note")
    printf 'NOTE: %s\n' "$note"
}

validation__command_text() {
    local rendered=""
    local part
    for part in "$@"; do
        printf -v part '%q' "$part"
        [[ -z "$rendered" ]] || rendered+=" "
        rendered+="$part"
    done
    printf '%s\n' "$rendered"
}

validation__slug() {
    local value="$1"
    value="${value,,}"
    value="${value//[^a-z0-9]/-}"
    while [[ "$value" == *--* ]]; do
        value="${value//--/-}"
    done
    value="${value#-}"
    value="${value%-}"
    printf '%.72s\n' "${value:-check}"
}

validation__extract_cause() {
    local log="$1"
    local cause

    cause="$(
        grep '^FINAL RESULT: FAIL — ' "$log" | tail -n 1 || true
    )"
    if [[ -n "$cause" ]]; then
        cause="${cause#FINAL RESULT: FAIL — }"
        cause="${cause#* — }"
        printf '%s\n' "$cause"
        return 0
    fi

    cause="$(
        awk '
            /^[[:space:]]*$/ { next }
            /^== .* ==$/ { next }
            /^PASS:/ { next }
            /^PASS checks:/ { next }
            /^FAIL checks:[[:space:]]*0[[:space:]]*$/ { next }
            /^Correctness result:/ { next }
            /^Resource observation:/ { next }
            /^Performance thresholds:/ { next }
            /^Phase .* validation (PASSED|FAILED)\.$/ { next }

            /command not found/ ||
            /No such file or directory/ ||
            /too many arguments in call/ ||
            /undefined:/ ||
            /syntax error/ ||
            /Traceback \(most recent call last\)/ ||
            /^fatal:/ ||
            /^ERROR:/ ||
            /^FAIL: / {
                print
                exit
            }
        ' "$log"
    )"

    if [[ -z "$cause" ]]; then
        cause="$(
            awk '
                /^[[:space:]]*$/ { next }
                { line=$0 }
                END { print line }
            ' "$log"
        )"
    fi

    printf '%s\n' "${cause:-command failed without diagnostic output}"
}

validation_run() {
    local name="$1"
    shift

    local sequence=$(( ${#VALIDATION_CHECK_NAMES[@]} + 1 ))
    local slug
    slug="$(validation__slug "$name")"
    local log="$VALIDATION_REPORT_DIR/$(printf '%03d' "$sequence")-$slug.log"
    local command_text
    command_text="$(validation__command_text "$@")"

    printf '\n== %s ==\n' "$name"
    printf 'Command: %s\n' "$command_text"
    printf 'Log: %s\n' "$log"

    local status
    if "$@" > >(tee "$log") 2>&1; then
        status=0
    else
        status=$?
    fi

    VALIDATION_CHECK_NAMES+=("$name")
    VALIDATION_CHECK_EXIT+=("$status")
    VALIDATION_CHECK_COMMAND+=("$command_text")
    VALIDATION_CHECK_LOG+=("$log")

    if (( status == 0 )); then
        VALIDATION_CHECK_STATUS+=("PASS")
        VALIDATION_CHECK_CAUSE+=("")
        VALIDATION_PASS_COUNT=$((VALIDATION_PASS_COUNT + 1))
        printf 'PASS: %s\n' "$name"
        return 0
    fi

    local cause
    cause="$(validation__extract_cause "$log")"
    VALIDATION_CHECK_STATUS+=("FAIL")
    VALIDATION_CHECK_CAUSE+=("$cause")
    VALIDATION_FAIL_COUNT=$((VALIDATION_FAIL_COUNT + 1))
    printf 'FAIL: %s\n' "$name"
    printf 'CAUSE: %s\n' "$cause"
    return "$status"
}

validation_skip() {
    local name="$1"
    local reason="$2"
    VALIDATION_SKIP_NAMES+=("$name")
    VALIDATION_SKIP_REASONS+=("$reason")
    VALIDATION_SKIP_COUNT=$((VALIDATION_SKIP_COUNT + 1))
    printf 'SKIP: %s — %s\n' "$name" "$reason"
}

validation__print_failure() {
    local index="$1"
    printf 'Check: %s\n' "${VALIDATION_CHECK_NAMES[$index]}"
    printf 'Exit status: %s\n' "${VALIDATION_CHECK_EXIT[$index]}"
    printf 'Command: %s\n' "${VALIDATION_CHECK_COMMAND[$index]}"
    printf 'Cause: %s\n' "${VALIDATION_CHECK_CAUSE[$index]}"
    printf 'Log: %s\n' "${VALIDATION_CHECK_LOG[$index]}"
}

validation_report_finish() {
    local report_path="$1"
    local latest_path="${2:-}"
    local temp_path="$VALIDATION_REPORT_DIR/.final-report.tmp"

    local first_failure=-1
    local index
    for index in "${!VALIDATION_CHECK_STATUS[@]}"; do
        if [[ "${VALIDATION_CHECK_STATUS[$index]}" == "FAIL" ]]; then
            first_failure="$index"
            break
        fi
    done

    {
        printf '\n============================================================\n'
        printf 'IRON ATLAS VALIDATION REPORT\n'
        printf 'Title: %s\n' "$VALIDATION_REPORT_TITLE"
        printf 'PASS checks: %d\n' "$VALIDATION_PASS_COUNT"
        printf 'FAIL checks: %d\n' "$VALIDATION_FAIL_COUNT"
        printf 'SKIPPED checks: %d\n' "$VALIDATION_SKIP_COUNT"

        if (( VALIDATION_FAIL_COUNT != 0 )); then
            printf '\nPRIMARY FAILURE\n'
            validation__print_failure "$first_failure"

            local primary_cause="${VALIDATION_CHECK_CAUSE[$first_failure]}"
            local -A seen_causes=()
            seen_causes["$primary_cause"]=1
            local additional_header=false
            local cascade_header=false

            for index in "${!VALIDATION_CHECK_STATUS[@]}"; do
                [[ "$index" == "$first_failure" ]] && continue
                [[ "${VALIDATION_CHECK_STATUS[$index]}" != "FAIL" ]] && continue

                local cause="${VALIDATION_CHECK_CAUSE[$index]}"
                if [[ "$cause" == "$primary_cause" ]]; then
                    if [[ "$cascade_header" == false ]]; then
                        printf '\nCASCADED FAILURES WITH THE SAME ROOT CAUSE\n'
                        cascade_header=true
                    fi
                    printf -- '- %s\n' "${VALIDATION_CHECK_NAMES[$index]}"
                    continue
                fi

                if [[ -n "${seen_causes[$cause]+x}" ]]; then
                    if [[ "$cascade_header" == false ]]; then
                        printf '\nCASCADED FAILURES WITH A REPEATED CAUSE\n'
                        cascade_header=true
                    fi
                    printf -- '- %s — %s\n' \
                        "${VALIDATION_CHECK_NAMES[$index]}" "$cause"
                    continue
                fi

                seen_causes["$cause"]=1
                if [[ "$additional_header" == false ]]; then
                    printf '\nADDITIONAL UNIQUE FAILURES\n'
                    additional_header=true
                fi
                printf '\n'
                validation__print_failure "$index"
            done
        fi

        if (( VALIDATION_SKIP_COUNT != 0 )); then
            printf '\nSKIPPED DEPENDENT CHECKS\n'
            for index in "${!VALIDATION_SKIP_NAMES[@]}"; do
                printf -- '- %s — %s\n' \
                    "${VALIDATION_SKIP_NAMES[$index]}" \
                    "${VALIDATION_SKIP_REASONS[$index]}"
            done
        fi

        if (( ${#VALIDATION_REPORT_NOTES[@]} != 0 )); then
            printf '\nREPORT NOTES\n'
            for index in "${!VALIDATION_REPORT_NOTES[@]}"; do
                printf -- '- %s\n' "${VALIDATION_REPORT_NOTES[$index]}"
            done
        fi

        printf '\nPer-check logs: %s\n' "$VALIDATION_REPORT_DIR"
        printf 'Final report: %s\n' "$report_path"
        printf '============================================================\n'

        if (( VALIDATION_FAIL_COUNT == 0 )); then
            printf 'FINAL RESULT: PASS — %s\n' "$VALIDATION_REPORT_TITLE"
        else
            printf 'FINAL RESULT: FAIL — %s — %s\n' \
                "${VALIDATION_CHECK_NAMES[$first_failure]}" \
                "${VALIDATION_CHECK_CAUSE[$first_failure]}"
        fi
    } | tee "$temp_path"

    mv "$temp_path" "$report_path"
    if [[ -n "$latest_path" ]]; then
        mkdir -p "$(dirname "$latest_path")"
        cp "$report_path" "$latest_path"
    fi

    (( VALIDATION_FAIL_COUNT == 0 ))
}
