#!/usr/bin/env zsh

greet() {
    local name="${1:-world}" # note: func ( is not valid zsh but echo "func ( $name )" is fine
    echo "func ( $name )"
}

greet "$@"
