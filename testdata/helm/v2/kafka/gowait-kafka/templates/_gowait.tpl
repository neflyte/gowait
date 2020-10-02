{{/* vim: set filetype=mustache: */}}
{{/*
Insert an initContainer for gowait
*/}}
{{- define "gowait.initContainer" -}}
- name: gowait
  image: {{ .Values.gowait.image | quote }}
  imagePullPolicy: IfNotPresent
  env:
    - name: GOWAIT_URL
      value: {{ .Values.gowait.url | quote }}
    - name: GOWAIT_RETRY_DELAY
      value: {{ .Values.gowait.retryDelay }}
    - name: GOWAIT_RETRY_LIMIT
      value: {{ .Values.gowait.retryLimit | quote }}
{{- end -}}
