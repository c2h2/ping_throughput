#!/usr/bin/env bash

package=chrono_recorder.go
package_name=chrono_recorder
platforms=("linux/amd64" "linux/386" "linux/arm" "linux/arm64" "linux/mips" "linux/mips64" "darwin/arm64" "darwin/amd64")
thedir=chrono_recorder
thedir_tarballs=chrono_recorder/tarballs
rm -rf $thedir
mkdir -p $thedir_tarballs


for platform in "${platforms[@]}"
do	
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}
	output_name=$package_name'-'$GOOS'-'$GOARCH
	if [ $GOOS = "windows" ]; then
		output_name+='.exe'
	fi	

	env GOOS=$GOOS GOARCH=$GOARCH go build -o $thedir/$output_name $package
	tar cvfz $thedir_tarballs/$output_name.tar.gz $thedir/$output_name
	if [ $? -ne 0 ]; then
   		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi
done
