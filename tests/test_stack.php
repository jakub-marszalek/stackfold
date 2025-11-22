<?php

$tests = [];

// Test 1: PHP Version
$tests['php_version'] = phpversion();

// Test 2: Extensions
$tests['extensions'] = [
    'pdo' => extension_loaded('pdo'),
    'pdo_pgsql' => extension_loaded('pdo_pgsql'),
    'redis' => extension_loaded('redis'),
    'intl' => extension_loaded('intl'),
];

// Test 3: PostgreSQL
try {
    $pdo = new PDO('pgsql:host=db;dbname=boilerplate', 'admin', 'Admin123');
    $tests['database'] = 'connected';
} catch (Exception $e) {
    $tests['database'] = 'error: ' . $e->getMessage();
}

// Test 4: Redis
try {
    $redis = new Redis();
    $redis->connect('redis');
    $redis->set('test_key', 'test_value');
    $value = $redis->get('test_key');
    $tests['redis'] = $value === 'test_value' ? 'connected' : 'error: value mismatch';
} catch (Exception $e) {
    $tests['redis'] = 'error: ' . $e->getMessage();
}

// Test 5: DNS Resolution
$tests['dns_tests'] = [
    'db' => gethostbyname('db'),
    'redis' => gethostbyname('redis'),
    'nginx' => gethostbyname('nginx'),
];

echo json_encode($tests, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES);