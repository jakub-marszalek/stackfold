<?php

/*
 * Copyright (c) 2025. Jakub Marszałek / hesoyamcc
 * All Rights Reserved.
 */

declare(strict_types=1);

namespace Health;

require_once __DIR__ . '/../../../vendor/autoload.php';

use Health\Stack\Checks\NetworkCheck;
use Health\Stack\Checks\PhpCheck;
use Health\Stack\Checks\PostgreCheck;
use Health\Stack\Checks\RedisCheck;
use Health\Stack\HealthChecker;

error_reporting(E_ALL);
ini_set('display_errors', 0);
ini_set('log_errors', 0);
ini_set('error_log', 'health_check_error.log');

// Create all checks with a health checker
$checks = [
    'php' => new PhpCheck(),
    'database' => new PostgreCheck(),
    'redis' => new RedisCheck(),
    'network' => new NetworkCheck(),
];

$healthChecker = new HealthChecker();

// Register all checks
foreach ($checks as $name => $check) {
    $healthChecker->register($name, $check);
}

// Execute checks
$result = $healthChecker->run();

// Return result as JSON
header('Content-Type: application/json; charset=utf-8');
echo json_encode($result, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES);
