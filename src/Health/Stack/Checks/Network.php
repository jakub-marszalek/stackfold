<?php

/*
 * Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
 * Licensed under the MIT License. See LICENSE file in the project root.
 */

declare(strict_types=1);

namespace Health\Stack\Checks;

use Health\Stack\Contracts\HealthCheckInterface;

class Network implements HealthCheckInterface
{
    private const array STACK_SERVICES = ['db', 'redis'];

    /**
     * @var array
     */
    private array $result = [
        'hostname' => null,
        'dns_resolution' => []
    ];

    /**
     * Check the network connectivity.
     *
     * @return array
     */
    public function check(): array
    {
        $this->result['hostname'] = gethostname();
        $this->result['dns_resolution'] = $this->resolveDns();
        return $this->result;
    }

    /**
     * Check if all checks were passed.
     *
     * @return bool
     */
    public function isHealthy(): bool
    {
        if (!$this->result['hostname']) {
            return false;
        }
        foreach ($this->result['dns_resolution'] as $ip) {
            if ($ip === null) {
                return false;
            }
        }
        return true;
    }

    /**
     * Resolve DNS for all listed services.
     *
     * @return array
     */
    private function resolveDns(): array
    {
        $resolved = [];
        foreach (self::STACK_SERVICES as $hostname) {
            $records = dns_get_record($hostname . '.', DNS_A | DNS_AAAA);
            if (empty($records)) {
                $resolved[$hostname] = null;
                continue;
            }
            $ipAddresses = $this->getIpAddresses($records);
            $resolved[$hostname] = $ipAddresses ?: null;
        }
        return $resolved;
    }

    /**
     * Get IP addresses from DNS records.
     *
     * @param array $records
     * @return array
     */
    private function getIpAddresses(array $records): array
    {
        $ipAddresses = [];
        foreach ($records as $record) {
            if (isset($record['ip'])) {
                $ipAddresses[] = $record['ip'];
            }
            if (isset($record['ipv6'])) {
                $ipAddresses[] = $record['ipv6'];
            }
        }
        return $ipAddresses;
    }
}
