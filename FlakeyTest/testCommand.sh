g=0; for n in $(gsutil ls gs://istio-circleci/master/test-integration-kubernetes/*/artifacts/junit.xml); do foo=$(cut -d "," -f 2 <<< $(cut -d ":" -f 2 <<< $(gsutil stat $n | sed -n 3p))); echo $n; gsutil cp $n "gs://istio-flakey-test/$data_folder/out-$foo-$g.xml"; ((++g)); done
