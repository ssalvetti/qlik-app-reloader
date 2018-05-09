$apps = @("firstApp","secondApp","thirdApp")
For ($i=0; $i -lt $apps.Length; $i++) {
    $app = $apps[$i]
    Start-Process -NoNewWindow -FilePath ".\qlik-app-reloader.exe" -ArgumentList "-app C:\\Users\\User\\Path\\$app"
}