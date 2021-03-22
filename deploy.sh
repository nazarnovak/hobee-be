go mod tidy &&
go mod vendor &&
cp -r vendor/* ../. &&
cp -r api ../github.com/nazarnovak/hobee-be &&
cp -r config ../github.com/nazarnovak/hobee-be &&
cp -r db ../github.com/nazarnovak/hobee-be &&
cp -r pkg ../github.com/nazarnovak/hobee-be &&
gcloud app deploy --stop-previous-version

printf "\nDon't forget to run ./cleanup after deploying, to remove old stopped instances. Save \$\$\$"
