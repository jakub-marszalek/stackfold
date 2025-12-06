<?php

/*
 * Copyright (c) 2025. Jakub Marszałek / hesoyamcc
 * All Rights Reserved.
 */

/**
 * Docker Stack Health Check
 * Tests all critical services and returns a JSON response
 */

error_reporting(E_ALL);
ini_set('display_errors', 1);

$results = [];
$allPassed = true;

// ==================== PHP ====================
$results['php'] = [
    'version' => phpversion(),
    'extensions' => [
        'pdo' => extension_loaded('pdo'),
        'pdo_pgsql' => extension_loaded('pdo_pgsql'),
        'redis' => extension_loaded('redis'),
        'intl' => extension_loaded('intl'),
        'json' => extension_loaded('json'),
    ],
    'memory_limit' => ini_get('memory_limit'),
    'max_execution_time' => ini_get('max_execution_time'),
];

// Check if necessary extensions are loaded
$requiredExtensions = ['pdo', 'pdo_pgsql', 'redis', 'intl'];
foreach ($requiredExtensions as $ext) {
    if (!$results['php']['extensions'][$ext]) {
        $allPassed = false;
    }
}

// ==================== PostgreSQL ====================
$results['database'] = [
    'status' => 'disconnected',
    'version' => 'unknown',
    'error' => null,
];

try {
    $pdo = new PDO(
        'pgsql:host=db;port=5432;dbname=boilerplate',
        'admin',
        'Admin123',
        [PDO::ATTR_TIMEOUT => 5]
    );

    $version = $pdo->query('SELECT version()')->fetchColumn();
    $results['database']['status'] = 'connected';
    $results['database']['version'] = $version;
    $results['database']['dsn'] = 'pgsql://db:5432/boilerplate';

} catch (PDOException $e) {
    $results['database']['error'] = $e->getMessage();
    $allPassed = false;
}

// ==================== Redis ====================
$results['cache'] = [
    'status' => 'disconnected',
    'version' => 'unknown',
    'error' => null,
    'test_value' => null,
];

try {
    $redis = new Redis();
    $redis->connect('redis', 6379, 5);

    // Test ping
    if ($redis->ping()) {
        $results['cache']['status'] = 'connected';

        // Test set/get
        $testKey = 'health_check_' . time();
        $testValue = 'test_' . uniqid();

        $redis->set($testKey, $testValue);
        $getValue = $redis->get($testKey);

        $results['cache']['test_value'] = ($getValue === $testValue) ? 'passed' : 'failed';
        $results['cache']['info'] = $redis->info('server');

        $redis->del($testKey);
    }
} catch (Exception $e) {
    $results['cache']['error'] = $e->getMessage();
    $allPassed = false;
}

// ==================== File System ====================
$results['filesystem'] = [
    'public_dir' => is_dir('public') ? 'exists' : 'missing',
    'storage_writable' => is_writable('.') ? 'yes' : 'no',
    'vendor_dir' => is_dir('vendor') ? 'exists' : 'missing',
    'node_modules' => is_dir('node_modules') ? 'exists' : 'missing',
];

// ==================== Network ====================
$results['network'] = [
    'hostname' => gethostname(),
    'dns_resolution' => [
        'db' => gethostbyname('db'),
        'redis' => gethostbyname('redis'),
    ],
];

// ==================== Environment ====================
$results['environment'] = [
    'app_env' => getenv('APP_ENV') ?: 'not set',
    'app_debug' => getenv('APP_DEBUG') ?: 'not set',
    'php_sapi' => php_sapi_name(),
    'os' => php_uname('s'),
];

// ==================== Final Status ====================
$results['summary'] = [
    'overall_status' => $allPassed ? 'PASSED' : 'FAILED',
    'timestamp' => date('Y-m-d H:i:s'),
    'critical_services' => [
        'php' => 'OK',
        'postgresql' => $results['database']['status'],
        'redis' => $results['cache']['status'],
        'nginx' => 'OK (running if you see this)',
        'node' => 'OK (check port 5173)',
    ],
];

// ==================== Output ====================
header('Content-Type: application/json; charset=utf-8');
echo json_encode($results, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES);
