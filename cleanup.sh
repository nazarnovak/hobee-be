stopped_versions=$(gcloud app versions list | grep 'STOPPED' | awk -F ' ' '{print $2}')

if [[ $stopped_versions == "" ]]; then
	echo 'No versions to remove';
	exit 0;
fi

gcloud app versions delete $stopped_versions
