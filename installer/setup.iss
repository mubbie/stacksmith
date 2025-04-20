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
PrivilegesRequired=admin
UserInfoPage=yes
ShowTasksTreeLines=yes
OutputDir=Output
OutputBaseFilename=stacksmith-setup
SetupIconFile=icon.ico
Compression=lzma
SolidCompression=yes

[Files]
Source: "..\artifact\stacksmith.exe"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\Stacksmith (CLI)"; Filename: "{app}\stacksmith.exe"
Name: "{group}\Uninstall Stacksmith"; Filename: "{uninstallexe}"

[Run]
Filename: "{app}\stacksmith.exe"; Description: "Run Stacksmith"; Flags: postinstall shellexec skipifsilent

[UninstallDelete]
Type: files; Name: "{app}\stacksmith.exe"

[Tasks]
Name: addtopath; Description: "Add Stacksmith to system PATH"; \
    GroupDescription: "Additional Options:"; Flags: unchecked

[Registry]
Root: HKLM; Subkey: "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"; \
    ValueType: expandsz; ValueName: "Path"; ValueData: "{app}"; Flags: preservestringtype uninsdeletevalue; Tasks: addtopath

[Code]
const
  SMTO_ABORTIFHUNG = 2;
  WM_SETTINGCHANGE = $1A;

function SendMessageTimeout(hWnd: HWND; Msg: LongWord; wParam, lParam: Longint; 
  fuFlags, uTimeout: LongWord; var lpdwResult: HWND): Integer;
  external 'SendMessageTimeoutA@user32.dll stdcall';
  
function FindWindowByClassName(lpClassName: string): HWND;
  external 'FindWindowA@user32.dll stdcall';
  
procedure RefreshEnvironment;
var
  Wnd: HWND;
begin
  Wnd := FindWindowByClassName('Progman');
  if Wnd <> 0 then
    SendMessageTimeout(Wnd, WM_SETTINGCHANGE, 0, 0, SMTO_ABORTIFHUNG, 5000, Wnd);
end;

procedure CurStepChanged(CurStep: TSetupStep);
begin
  if CurStep = ssPostInstall then
  begin
    if WizardIsTaskSelected('addtopath') then
      RefreshEnvironment;
  end;
end;
