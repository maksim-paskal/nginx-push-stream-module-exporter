#!/bin/sh

helm lint --strict ./comet-exporter
helm template ./comet-exporter | kubectl apply --dry-run --validate -f -