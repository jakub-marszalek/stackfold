<?php

echo json_encode([
    'status' => 'ok',
    'php_version' => phpversion(),
    'extensions' => [
        'pdo' => extension_loaded('pdo'),
        'pdo_pgsql' => extension_loaded('pdo_pgsql'),
        'redis' => extension_loaded('redis'),
        'intl' => extension_loaded('intl'),
    ]
]);
