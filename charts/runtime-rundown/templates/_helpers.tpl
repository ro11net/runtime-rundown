{{/*
Expand the name of the chart.
*/}}
{{- define "flux-training-app.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "flux-training-app.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "flux-training-app.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "flux-training-app.labels" -}}
helm.sh/chart: {{ include "flux-training-app.chart" . }}
{{ include "flux-training-app.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "flux-training-app.selectorLabels" -}}
app.kubernetes.io/name: {{ include "flux-training-app.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "flux-training-app.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "flux-training-app.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create a default fully qualified postgresql name.
*/}}
{{- define "flux-training-app.postgresql.fullname" -}}
{{- $name := default "postgresql" .Values.postgresql.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Get the PostgreSQL secret name
*/}}
{{- define "flux-training-app.postgresql.secretName" -}}
{{- if .Values.postgresql.auth.existingSecret -}}
    {{- .Values.postgresql.auth.existingSecret -}}
{{- else -}}
    {{- include "flux-training-app.postgresql.fullname" . -}}
{{- end -}}
{{- end -}}
