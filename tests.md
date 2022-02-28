go build -o translate && ./translate -path ./content/en/authors/ken-moini/_index.md -dest ./content/es/authors/ken-moini/_index.md -from en -to es -recursive

go build -o translate && ./translate -path ./content/en/post/2022/01/windows-vr-vm-on-libvirt.md -dest ./content/es/post/2022/01/windows-vr-vm-on-libvirt.md -from en -to es -recursive

go build -o translate && ./translate -path ./content/en/authors/ -dest ./content/es/authors/ -from en -to es -recursive