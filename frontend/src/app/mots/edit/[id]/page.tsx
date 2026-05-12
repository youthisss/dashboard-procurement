"use client";

import { SimpleMasterEdit } from "@/components/master-data/simple-master-pages";

export default function MotEdit() {
  return (
    <SimpleMasterEdit
      title="Edit MOT"
      fieldLabel="MOT Name"
      placeholder="e.g. Truck"
      requiredMessage="MOT name is required"
    />
  );
}
