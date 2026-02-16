<?php

declare(strict_types=1);

namespace Health;

require_once __DIR__ . '/../../../vendor/autoload.php';

use Health\Stack\Checks\Network;
use Health\Stack\Checks\Php;
use Health\Stack\Checks\Postgres;
use Health\Stack\Checks\RedisCheck;
use Health\Stack\HealthChecker;

// Configure ini settings
error_reporting(E_ALL);
ini_set('display_errors', 0);
ini_set('log_errors', 0);
ini_set('error_log', 'health_check_error.log');

// Create all checks and the health checker
$checks = [
    'php' => new Php(),
    'database' => new Postgres(),
    'redis' => new RedisCheck(),
    'network' => new Network(),
];
$healthChecker = new HealthChecker();

// Register all checks and execute them
foreach ($checks as $name => $check) {
    $healthChecker->register($name, $check);
}
$result = $healthChecker->run();

// Return the results in JSON format
header('Content-Type: application/json; charset=utf-8');
echo json_encode($result, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES);
