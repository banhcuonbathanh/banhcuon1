"use client";

import React, { useState } from "react";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { ChevronDown, ChevronRight } from "lucide-react";

import { loggerPaths, LogPath } from "@/lib/logger/loggerConfig";

// Utility functions
const getLogRange = (ids: number[]): string => {
  const min = Math.min(...ids);
  const max = Math.max(...ids);
  return `${min}-${max}`;
};

const getComponentName = (path: string): string => {
  const cleanPath = path.replace(/\.[jt]sx?$/, "");
  const fileName = cleanPath.split("/").pop() || "";
  const parts = fileName.split("-");
  return parts
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
};

interface LogDetailsProps {
  config: LogPath;
}

const LogDetails: React.FC<LogDetailsProps> = ({ config }) => {
  return (
    <div className="p-4 space-y-4">
      <div className="grid grid-cols-3 gap-4">
        <div className=" p-3 rounded-lg">
          <h4 className="font-semibold mb-2 text-sm">Enabled Log IDs</h4>
          <div className="text-sm text-gray-600">
            {config.enabledLogIds.join(", ")}
          </div>
        </div>
        <div className=" p-3 rounded-lg">
          <h4 className="font-semibold mb-2 text-sm">Disabled Log IDs</h4>
          <div className="text-sm text-gray-600">
            {config.disabledLogIds.length
              ? config.disabledLogIds.join(", ")
              : "None"}
          </div>
        </div>
        <div className=" p-3 rounded-lg">
          <h4 className="font-semibold mb-2 text-sm">Total Log IDs</h4>
          <div className="text-sm text-gray-600">
            {config.logIds.join(", ")}
          </div>
        </div>
      </div>

      <div className="mt-4">
        <h4 className="font-semibold mb-3">Log Details</h4>
        <div className="grid grid-cols-1 gap-2 max-h-96 overflow-y-auto">
          {Object.entries(config.logDescriptions).map(([id, log]) => (
            <div key={id} className=" p-3 rounded-lg">
              <div className="flex items-center justify-between mb-2">
                <span className="font-medium">ID: {id}</span>
                <Badge
                  variant={log.status === "enabled" ? "default" : "secondary"}
                >
                  {log.status}
                </Badge>
              </div>
              <p className="text-sm text-gray-600">{log.description}</p>
              <p className="text-xs text-gray-500 mt-1">
                Location: {log.location}
              </p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

const LoggerSummary = () => {
  const [expandedPaths, setExpandedPaths] = useState<string[]>([]);
  const [paths, setPaths] = useState<LogPath[]>(loggerPaths);

  const togglePath = (path: string) => {
    setExpandedPaths((prev) =>
      prev.includes(path) ? prev.filter((p) => p !== path) : [...prev, path]
    );
  };

  const toggleStatus = (pathToToggle: string) => {
    setPaths((prevPaths) =>
      prevPaths.map((config) =>
        config.path === pathToToggle
          ? { ...config, enabled: !config.enabled }
          : config
      )
    );
  };

  return (
    <div className="p-6 max-w-6xl mx-auto space-y-6">
      <h1 className="text-3xl font-bold">Logger Configuration Summary</h1>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-48">Component</TableHead>
              <TableHead className="w-32">Log IDs</TableHead>
              <TableHead className="w-32">Status</TableHead>
              <TableHead>Path</TableHead>
              <TableHead className="w-20 text-right">Details</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {paths.map((config: LogPath) => (
              <React.Fragment key={config.path}>
                <TableRow>
                  <TableCell className="font-medium">
                    {getComponentName(config.description)}
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="ghost"
                      onClick={() => togglePath(config.path)}
                      className="px-2 py-1 h-auto"
                    >
                      {getLogRange(config.enabledLogIds)}
                      {expandedPaths.includes(config.path) ? (
                        <ChevronDown className="h-4 w-4 ml-1" />
                      ) : (
                        <ChevronRight className="h-4 w-4 ml-1" />
                      )}
                    </Button>
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="ghost"
                      onClick={() => toggleStatus(config.path)}
                      className="px-2 py-1 h-auto"
                    >
                      <Badge
                        variant={config.enabled ? "default" : "secondary"}
                        className="cursor-pointer"
                      >
                        {config.enabled ? "Enabled" : "Disabled"}
                      </Badge>
                    </Button>
                  </TableCell>
                  <TableCell className="font-mono text-sm text-gray-600">
                    {config.path}
                  </TableCell>
                  <TableCell className="text-right">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => togglePath(config.path)}
                    >
                      {expandedPaths.includes(config.path) ? (
                        <ChevronDown className="h-4 w-4" />
                      ) : (
                        <ChevronRight className="h-4 w-4" />
                      )}
                    </Button>
                  </TableCell>
                </TableRow>
                {expandedPaths.includes(config.path) && (
                  <TableRow>
                    <TableCell colSpan={5} className="">
                      <LogDetails config={config} />
                    </TableCell>
                  </TableRow>
                )}
              </React.Fragment>
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  );
};

export default LoggerSummary;
