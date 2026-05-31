$conns = Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue
foreach ($c in $conns) {
    Stop-Process -Id $c.OwningProcess -Force -ErrorAction SilentlyContinue
    Write-Host "Killed PID $($c.OwningProcess)"
}
Write-Host "Done"
