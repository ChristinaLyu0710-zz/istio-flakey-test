gcloud auth application-default login; g=0; for n in $(gsutil ls {gs://istio-circleci/master/*/*/artifacts/junit.xml,gs://istio-prow/logs/*master/*/artifacts/junit.xml}); do foo=$(cut -d "," -f 2 <<< $(cut -d ":" -f 2 <<< $(gsutil stat $n | sed -n 3p))); gsutil cp $n "gs://istio-flakey-test/temp/out-$foo-$g.xml"; ((++g)); done