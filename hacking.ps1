$env:GOOS = "linux"
rm bootstrap; 
rm bootstrap.zip;
go build;
mv lambda bootstrap;
Compress-Archive bootstrap bootstrap.zip;
rm bootstrap;