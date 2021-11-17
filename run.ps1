param([Int32]$n=3)
#Run with -nodes <number of nodes>
#go run . -name node1 -address 127.0.0.1 -sport 8080 -bport 8081 -cport 8081
#go run . -name node2 -address 127.0.0.1 -sport 8082 -bport 8083 -cport 8081

$names = @('node0', 'node1', 'node2', 'node3', 'node4')
$sports = @('8080', '8081', '8082', '8083', '8084')
$ips = ''
$delays = @(5, 5, 5, 5, 5)

for ($i = 0; $i -lt $n; $i++) {
    $ips = ''
    if ($i -lt $names.count) {
        $name = $names[$i]
        $sport = $sports[$i]
        $delay = $delays[$i]
        
        for ($j = 0; $j -lt $n; $j++) {
            
            if ($n -ne 1) {
                if ($sports[$j] -eq $sport) {
                    continue
                }
            }

            if ($ips -eq '') {
                $ips += $sports[$j]
            } else {
                $ips += ',' + $sports[$j]
            }
        }
    } else {
        "WARNING: Max number of nodes for this script has been created (" + $names.count + "). => Not creating any more nodes!"
        break
    }

    "STARTING GO => $name | Server port: $sport | IPs: $ips | Delay: $delay |"

    $Command = 'cmd /c start powershell -NoExit -Command {
            $host.UI.RawUI.WindowTitle = "Node - ' + $name + '";
            $host.UI.RawUI.BackgroundColor = "black";
            $host.UI.RawUI.ForegroundColor = "white";
            Clear-Host;
            cd node;
            go run . -name ' + $name + ' -sport ' + $sport + ' -ips ' + $ips + ' -delay ' + $delay + ';
        }'

    invoke-expression -Command $Command
}