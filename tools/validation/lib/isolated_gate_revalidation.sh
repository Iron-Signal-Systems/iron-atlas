#!/usr/bin/env bash
set -Eeuo pipefail

isolated_gate_revalidate() {
    if (( $# != 3 )); then
        printf 'isolated_gate_revalidate requires: repository commit validator-path\n' >&2
        return 2
    fi

    local source_repo="$1"
    local commit="$2"
    local validator_path="$3"
    local tmp rc=0

    [[ "$validator_path" != /* ]] || {
        printf 'validator path must be repository-relative\n' >&2
        return 2
    }
    [[ "$validator_path" != *".."* ]] || {
        printf 'validator path must not contain parent traversal\n' >&2
        return 2
    }

    tmp="$(mktemp -d)" || return 1

    git clone --quiet --shared --no-checkout "$source_repo" "$tmp/repo" || rc=$?
    if (( rc == 0 )); then
        git -C "$tmp/repo" switch --quiet -C dev "$commit" || rc=$?
    fi
    if (( rc == 0 )); then
        if [[ ! -x "$tmp/repo/$validator_path" ]]; then
            printf 'isolated validator is missing or not executable: %s\n' "$validator_path" >&2
            rc=1
        else
            "$tmp/repo/$validator_path"
            rc=$?
        fi
    fi

    rm -rf "$tmp" || true
    return "$rc"
}

isolated_python_validator_revalidate() {
    if (( $# != 3 )); then
        printf 'isolated_python_validator_revalidate requires: repository commit validator-path\n' >&2
        return 2
    fi

    local source_repo="$1"
    local commit="$2"
    local validator_path="$3"
    local tmp rc=0

    [[ "$validator_path" != /* ]] || {
        printf 'validator path must be repository-relative\n' >&2
        return 2
    }
    [[ "$validator_path" != *".."* ]] || {
        printf 'validator path must not contain parent traversal\n' >&2
        return 2
    }
    [[ "$validator_path" == *.py ]] || {
        printf 'isolated Python validator path must end in .py\n' >&2
        return 2
    }

    tmp="$(mktemp -d)" || return 1

    git clone --quiet --shared --no-checkout "$source_repo" "$tmp/repo" || rc=$?
    if (( rc == 0 )); then
        git -C "$tmp/repo" switch --quiet -C dev "$commit" || rc=$?
    fi
    if (( rc == 0 )); then
        if [[ ! -f "$tmp/repo/$validator_path" ]]; then
            printf 'isolated Python validator is missing: %s\n' "$validator_path" >&2
            rc=1
        else
            (
                cd "$tmp/repo"
                python3 "$validator_path"
            )
            rc=$?
        fi
    fi

    rm -rf "$tmp" || true
    return "$rc"
}
