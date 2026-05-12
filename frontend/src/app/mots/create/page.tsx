"use client";

import { SimpleMasterCreate } from "@/components/master-data/simple-master-pages";

export default function MotCreate() {
  return (
    <SimpleMasterCreate
      title="Create MOT"
      fieldLabel="MOT Name"
      placeholder="e.g. Truck"
      requiredMessage="MOT name is required"
    />
  );
}
