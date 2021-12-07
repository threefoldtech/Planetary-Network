#!/bin/sh

if [ ! -d wixbin ];
then
  curl -LO https://github.com/wixtoolset/wix3/releases/download/wix3112rtm/wix311-binaries.zip
  if [ `md5sum wix311-binaries.zip | cut -f 1 -d " "` != "47a506f8ab6666ee3cc502fb07d0ee2a" ];
  then
    echo "wix package didn't match expected checksum"
    exit 1
  fi
  mkdir -p wixbin
  unzip -o wix311-binaries.zip -d wixbin || (
    echo "failed to unzip WiX"
    exit 1
  )
fi

GOOS=windows GOARCH=amd64 CGO_ENABLED=0 ./build.sh

PKGNAME="ThreeFoldNetworkConnector"
PKGVERSION="0.0.0.1"
PKGVERSIONMS="0.0.0.1"
[ "x64" == "x64" ] && \
  PKGGUID="77757838-1a23-40a5-a720-c3b43e0260cc" PKGINSTFOLDER="ProgramFiles64Folder" || \
  PKGGUID="54a3294e-a441-4322-aefb-3bb40dd022bb" PKGINSTFOLDER="ProgramFilesFolder"

if [ ! -d wintun ];
then
  curl -o wintun.zip https://www.wintun.net/builds/wintun-0.11.zip
  unzip wintun.zip
fi

PKGWINTUNDLL=wintun/bin/amd64/wintun.dll

PKGDISPLAYNAME="ThreeFoldNetworkConnector"

# Generate the wix.xml file
cat > wix.xml << EOF
<?xml version="1.0" encoding="windows-1252"?>
<Wix xmlns="http://schemas.microsoft.com/wix/2006/wi">
  <Product
    Name="${PKGDISPLAYNAME}"
    Id="*"
    UpgradeCode="${PKGGUID}"
    Language="1033"
    Codepage="1252"
    Version="${PKGVERSIONMS}"
    Manufacturer="github.com/threefoldtech">

    <Package
      Id="*"
      Keywords="Installer"
      Description="ThreeFoldNetworkConnector"
      Comments="ThreeFoldNetworkConnector"
      Manufacturer="github.com/threefoldtech"
      InstallerVersion="200"
      InstallScope="perMachine"
      Languages="1033"
      Compressed="yes"
      SummaryCodepage="1252" />

    <MajorUpgrade
      AllowDowngrades="yes" />

    <Media
      Id="1"
      Cabinet="Media.cab"
      EmbedCab="yes"
      CompressionLevel="high" />

    <Directory Id="TARGETDIR" Name="SourceDir">

      <Directory Id="DesktopFolder" Name="Desktop">
        <Component Id="ApplicationShortcutDesktop" Guid="c5119291-2aa3-4962-864a-9759c87beb63">
            <Shortcut Id="ApplicationDesktopShortcut"
                Name="ThreeFold Planetary Network"
                Description="Connects to the ThreeFold Planetary Network."
                Target="[ThreeFoldInstallFolder]ThreeFoldPlanetaryNetwork.exe"
                WorkingDirectory="ThreeFoldInstallFolder"/>
            <RemoveFolder Id="DesktopFolder" On="uninstall"/>
            <RegistryValue
                Root="HKCU"
                Key="Software/ThreeFoldPlanetaryNetwork"
                Name="installed"
                Type="integer"
                Value="1"
                KeyPath="yes"/>
        </Component>
      </Directory>


      <Directory Id="${PKGINSTFOLDER}" Name="PFiles">
        <Directory Id="ThreeFoldInstallFolder" Name="Threefold">

          <Component Id="MainExecutable" Guid="c5119291-2aa3-4962-864a-9759c87beb64">
            <File
              Id="ThreeFoldPlanetaryNetwork"
              Name="ThreeFoldPlanetaryNetwork.exe"
              DiskId="1"
              Source="ThreeFoldPlanetaryNetwork.exe"
              KeyPath="yes" />

            <File
              Id="Wintun"
              Name="wintun.dll"
              DiskId="1"
              Source="${PKGWINTUNDLL}" />
          </Component>

        </Directory>
      </Directory>
    </Directory>

    <Feature Id="ThreeFoldFeature" Title="ThreeFoldPlanetaryNetwork" Level="1">
      <ComponentRef Id="MainExecutable" />
      <ComponentRef Id="ApplicationShortcutDesktop" />
    </Feature>

    <InstallExecuteSequence>
    </InstallExecuteSequence>

  </Product>
</Wix>
EOF

# Generate the MSI
CANDLEFLAGS="-nologo"
LIGHTFLAGS="-nologo -spdb -sice:ICE71 -sice:ICE61"
wixbin/candle $CANDLEFLAGS -out ${PKGNAME}-${PKGVERSION}-x64.wixobj -arch x64 wix.xml && \
wixbin/light $LIGHTFLAGS -ext WixUtilExtension.dll -out ${PKGNAME}-${PKGVERSION}-x64.msi ${PKGNAME}-${PKGVERSION}-x64.wixobj
