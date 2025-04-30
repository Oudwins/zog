// // MyTable.tsx
// import React from "react";
// // import styles from "./MyTable.module.css"; // Optional: CSS module for styling

// // Define the type for a single row in the table
// interface TableRow {
//   schemaType: string;
//   data: string;
//   dest: string;
//   requiredError: string;
//   coercionError: string;
// }

// const MyTable: React.FC = () => {
//   // Define the table data array with the specified type
//   const tableData: TableRow[] = [
//     {
//       schemaType: "Bool()",
//       data: "true",
//       dest: "true",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: "false",
//       dest: "false",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: "nil",
//       dest: "false",
//       requiredError: "yes",
//       coercionError: "yes",
//     },
//     {
//       schemaType: "Bool()",
//       data: '""',
//       dest: "false",
//       requiredError: "yes",
//       coercionError: "yes",
//     },
//     {
//       schemaType: "Bool()",
//       data: '" "',
//       dest: "false",
//       requiredError: "yes",
//       coercionError: "yes",
//     },
//     {
//       schemaType: "Bool()",
//       data: '"on"',
//       dest: "true",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: '"off"',
//       dest: "false",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: '"true", "t", "T", "True", "TRUE"',
//       dest: "true",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: '"false", "f", "F", "FALSE", "False"',
//       dest: "false",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: '"test"',
//       dest: "false",
//       requiredError: "no",
//       coercionError: "yes",
//     },
//     {
//       schemaType: "Bool()",
//       data: "1",
//       dest: "true",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: "0",
//       dest: "false",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: "123",
//       dest: "false",
//       requiredError: "no",
//       coercionError: "yes",
//     },
//     {
//       schemaType: "String()",
//       data: '""',
//       dest: '""',
//       requiredError: "yes",
//       coercionError: "no",
//     },
//     {
//       schemaType: "String()",
//       data: '" "',
//       dest: '""',
//       requiredError: "yes",
//       coercionError: "no",
//     },
//     {
//       schemaType: "String()",
//       data: "nil",
//       dest: '""',
//       requiredError: "yes",
//       coercionError: "yes",
//     },
//     {
//       schemaType: "String()",
//       data: "any value",
//       dest: 'fmt.Sprintf("%v", value)',
//       requiredError: "no",
//       coercionError: "no",
//     },
//     // Add the rest of your data rows here
//   ];

//   return (
//     <table>
//       <thead>
//         <tr>
//           <th>Schema Type</th>
//           <th>Data</th>
//           <th>Dest</th>
//           <th>Required Error (Zero Value)</th>
//           <th>Coercion Error</th>
//         </tr>
//       </thead>
//       <tbody>
//         {tableData.map((row, index) => (
//           <tr key={index}>
//             <td>{row.schemaType}</td>
//             <td>{row.data}</td>
//             <td>{row.dest}</td>
//             <td>{row.requiredError}</td>
//             <td>{row.coercionError}</td>
//           </tr>
//         ))}
//       </tbody>
//     </table>
//   );
// };

// export default MyTable;

const schemaData = [
  {
    schemaType: "Bool()",
    data: "true",
    dest: "true",
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Bool()",
    data: "false",
    dest: "false",
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Bool()",
    data: "nil",
    dest: "false",
    requiredError: "yes",
    coercionError: "yes",
  },
  {
    schemaType: "Bool()",
    data: '""',
    dest: "false",
    requiredError: "no",
    coercionError: "yes",
  },
  {
    schemaType: "Bool()",
    data: '" "',
    dest: "false",
    requiredError: "no",
    coercionError: "yes",
  },
  {
    schemaType: "Bool()",
    data: "on",
    dest: "true",
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Bool()",
    data: "off",
    dest: "false",
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Bool()",
    data: '["true", "t", "T", "True", "TRUE"]',
    dest: "true",
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Bool()",
    data: '["false", "f", "F", "FALSE", "False"]',
    dest: "false",
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Bool()",
    data: "test",
    dest: "false",
    requiredError: "no",
    coercionError: "yes",
  },
  {
    schemaType: "Bool()",
    data: "1",
    dest: "true",
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Bool()",
    data: "0",
    dest: "false",
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Bool()",
    data: "123",
    dest: "false",
    requiredError: "no",
    coercionError: "yes",
  },
  {
    schemaType: "String()",
    data: '""',
    dest: '""',
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "String()",
    data: '" "',
    dest: '""',
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "String()",
    data: "nil",
    dest: '""',
    requiredError: "yes",
    coercionError: "yes",
  },
  {
    schemaType: "String()",
    data: "any value",
    dest: 'fmt.Sprintf("%v", value)',
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Int()",
    data: 0,
    dest: 0,
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Int()",
    data: 10,
    dest: 10,
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Int()",
    data: "nil",
    dest: 0,
    requiredError: "yes",
    coercionError: "yes",
  },
  {
    schemaType: "Int()",
    data: '""',
    dest: 0,
    requiredError: "no",
    coercionError: "yes",
  },
  {
    schemaType: "Int()",
    data: '" "',
    dest: 0,
    requiredError: "no",
    coercionError: "yes",
  },
  {
    schemaType: "Int()",
    data: "any string",
    dest: "strconv.Atoi(str)",
    requiredError: "no",
    coercionError: "depends",
  },
  {
    schemaType: "Int()",
    data: 6.29,
    dest: 6,
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Int()",
    data: "true",
    dest: 1,
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Int()",
    data: "false",
    dest: 0,
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Float()",
    data: 1.21,
    dest: 1.21,
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Float()",
    data: 0.0,
    dest: 0.0,
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Float()",
    data: "nil",
    dest: 0,
    requiredError: "yes",
    coercionError: "yes",
  },
  {
    schemaType: "Float()",
    data: '""',
    dest: 0,
    requiredError: "no",
    coercionError: "yes",
  },
  {
    schemaType: "Float()",
    data: '" "',
    dest: 0,
    requiredError: "no",
    coercionError: "yes",
  },
  {
    schemaType: "Float()",
    data: "any string",
    dest: "strconv.ParseFloat(str)",
    requiredError: "no",
    coercionError: "depends",
  },
  {
    schemaType: "Float()",
    data: "1",
    dest: 1.0,
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Time()",
    data: "time.Time{}",
    dest: "time.Time{}",
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Time()",
    data: "time.Now()",
    dest: "time.Now()",
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Time()",
    data: "nil",
    dest: "time.Time{}",
    requiredError: "yes",
    coercionError: "yes",
  },
  {
    schemaType: "Time()",
    data: '""',
    dest: "time.Time{}",
    requiredError: "no",
    coercionError: "yes",
  },
  {
    schemaType: "Time()",
    data: '" "',
    dest: "time.Time{}",
    requiredError: "no",
    coercionError: "yes",
  },
  {
    schemaType: "Time()",
    data: "unix_timestamp_ms",
    dest: "time.Unix(unix, 0)",
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Time()",
    data: "any string",
    dest: "time.Parse(format, str)",
    requiredError: "no",
    coercionError: "depends",
  },
  {
    schemaType: "Slice()",
    data: "[1]",
    dest: "[1]",
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Slice()",
    data: "[]",
    dest: "[]",
    requiredError: "no",
    coercionError: "no",
  },
  {
    schemaType: "Slice()",
    data: "nil",
    dest: "[null]",
    requiredError: "yes",
    coercionError: "no (error will show in the appropriate schema if any)",
  },
  {
    schemaType: "Slice()",
    data: '""',
    dest: '[""]',
    requiredError: "no",
    coercionError: "no (error will show in the appropriate schema if any)",
  },
  {
    schemaType: "Slice()",
    data: '" "',
    dest: '[" "]',
    requiredError: "no",
    coercionError: "no (error will show in the appropriate schema if any)",
  },
  {
    schemaType: "Slice()",
    data: "any_value",
    dest: "[value]",
    requiredError: "depends",
    coercionError: "no (error will show in the appropriate schema if any)",
  },
];

export default function MyTable() {
  return (
    <div className="">
      <table>
        <thead>
          <tr>
            <th>Schema Type</th>
            <th>Data</th>
            <th>Dest</th>
            <th>Required Error (Zero Value)</th>
            <th>Coercion Error</th>
          </tr>
        </thead>
        <tbody>
          {schemaData.map((row, index) => (
            <tr key={index}>
              <td>{row.schemaType}</td>
              <td>{row.data}</td>
              <td>{row.dest}</td>
              <td>{row.requiredError}</td>
              <td>{row.coercionError}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
