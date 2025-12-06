<?php

/*
 * Copyright (c) 2025. Jakub Marszałek / hesoyamcc
 * All Rights Reserved.
 */

declare(strict_types=1);

namespace Health\Stack\Checks;

use Health\Stack\Contracts\HealthCheckInterface;

class NetworkCheck implements HealthCheckInterface
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
     * Check network connectivity
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
     * Check if all checks passed
     *
     * @return bool
     */
    public function isHealthy(): bool
    {
        if (!$this->result['hostname']) {
            return false;
        }
        foreach ($this->result['dns_resolution'] as $ip) {
            if (!$ip) {
                return false;
            }
        }
        return true;
    }

    /**
     * Resolve DNS for listed services
     *
     * @return array
     */
    private function resolveDns(): array
    {
        $hostnames = self::STACK_SERVICES;
        $resolved = [];

        foreach ($hostnames as $hostname) {
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
     * Get IP addresses from DNS records
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
