; Stacksmith Installer Script
; Requires Inno Setup 6+
[Setup]
AppName=Stacksmith
AppVersion={#MyAppVersion}
AppPublisher=Mubbie
AppPublisherURL=https://github.com/mubbie/stacksmith
AppSupportURL=https://github.com/mubbie/stacksmith/issues
AppUpdatesURL=https://github.com/mubbie/stacksmith/releases
DefaultDirName={autopf}\Stacksmith
DisableProgramGroupPage=yes
OutputDir=Output
OutputBaseFilename=stacksmith-setup
SetupIconFile=installer/icon.ico
Compression=lzma
SolidCompression=yes

[Files]
Source: "artifact\stacksmith.exe"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\Stacksmith (CLI)"; Filename: "{app}\stacksmith.exe"
Name: "{group}\Uninstall Stacksmith"; Filename: "{uninstallexe}"

[Run]
Filename: "{app}\stacksmith.exe"; Description: "Run Stacksmith"; Flags: postinstall shellexec skipifsilent

[UninstallDelete]
Type: files; Name: "{app}\stacksmith.exe"