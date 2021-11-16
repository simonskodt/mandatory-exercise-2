param([Int32]$nodes=3)
#Run with -nodes <number of nodes>
#go run . -name node1 -address 127.0.0.1 -sport 8080 -bport 8081 -cport 8081
#go run . -name node2 -address 127.0.0.1 -sport 8082 -bport 8083 -cport 8081

$names = @('node0', 'node1', 'node2', 'node3')
$sports = @(8080, 8081, 8082, 8083)

for ($i = 0; $i -lt $nodes; $i++) {
    $name = $i

    if ($i -lt $names.count) {
        $name = $names[$i]
        $sport = $sports[$i]
    } else {
        "WARNING: Max number of nodes for this script has been created."
        "Not creating any more nodes!"
        break
    }

    $Command = 'cmd /c start powershell -NoExit -Command {
            $host.UI.RawUI.WindowTitle = "Node - ' + $name + '";
            $host.UI.RawUI.BackgroundColor = "black";
            $host.UI.RawUI.ForegroundColor = "white";
            Clear-Host;
            cd node;
            go run . -name ' + $name +' -address ' + $address + ' -sport ' + $sport + ' -bport ' + ';
        }'

    invoke-expression -Command $Command
}