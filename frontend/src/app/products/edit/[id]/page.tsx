"use client";

import { SimpleMasterEdit } from "@/components/master-data/simple-master-pages";

export default function ProductEdit() {
  return (
    <SimpleMasterEdit
      title="Edit Product"
      fieldLabel="Product Name"
      placeholder="e.g. CPO"
      requiredMessage="Product name is required"
    />
  );
}
