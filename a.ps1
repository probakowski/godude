$VSINSTALLDIR = $(vswhere.exe -latest -requires Microsoft.VisualStudio.Component.VC.Llvm.Clang -property installationPath)
$VCINSTALLDIR = Join-Path $VSINSTALLDIR "VC"
$LLVM_ROOT = Join-Path $VCINSTALLDIR "Tools\Llvm\x64"
echo $LLVM_ROOT
ls $LLVM_ROOT\bin