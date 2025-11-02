{{/*
Expand the name of the chart.
*/}}
{{- define "runtime-rundown.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "runtime-rundown.fullname" -}}
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
{{- define "runtime-rundown.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "runtime-rundown.labels" -}}
helm.sh/chart: {{ include "runtime-rundown.chart" . }}
{{ include "runtime-rundown.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "runtime-rundown.selectorLabels" -}}
app.kubernetes.io/name: {{ include "runtime-rundown.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "runtime-rundown.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "runtime-rundown.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create a default fully qualified postgresql name.
*/}}
{{- define "runtime-rundown.postgresql.fullname" -}}
{{- $name := default "postgresql" .Values.postgresql.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Get the PostgreSQL secret name
*/}}
{{- define "runtime-rundown.postgresql.secretName" -}}
{{- if .Values.postgresql.auth.existingSecret -}}
    {{- .Values.postgresql.auth.existingSecret -}}
{{- else -}}
    {{- include "runtime-rundown.postgresql.fullname" . -}}
{{- end -}}
{{- end -}}
