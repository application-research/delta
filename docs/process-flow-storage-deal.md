# Storage Deal Making Process flow

Delta is a deal-making service that enables users to make deals with Storage Providers. It is an application that allows users to upload files to the Filecoin network and get them stored by Storage Providers.

A content goes thru different stages in delta. 

## Deal Preparation
- To start, the user will initiate the process by uploading the desired content onto Delta's platform. Once uploaded, Delta will assign the content to a light node and create a corresponding record in its database. To ensure the security and reliability of the content, Delta will assign a miner and wallet to the content record, and create a piece commitment record in its database. This involves computing the piece commitment of the content, which essentially verifies that the content has not been altered in any way.
- Delta will create a deal proposal parameters record in its database, which outlines the terms and conditions of the proposed deal. Based on this record, Delta will then create a deal proposal record, which will serve as a formal agreement between the user and Delta regarding the storage and access of the content on Delta's platform.
- Once all the meta is available, Delta will then dispatch a job that will call a method called `makeStorageDeal` to make a storage deal proposal for the content.

## Making a Storage Deal
- Once all the meta is available, Delta will then dispatch a job that will call a method called `makeStorageDeal` to make a storage deal proposal for the content.
- The purpose of this method is to make a storage deal proposal for a given Content object and PieceCommitment object.
- The first thing this process does is to update the status of the Content object to indicate that a deal proposal is being made. If there is any error during this process, the content status is updated to indicate that the proposal has failed and the error is returned.
- Next, it retrieves the miner address and Filecoin client associated with the content, as well as the deal proposal for the content. If there is any error, the content status is updated to indicate that the proposal has failed and the error is returned.
- It then sets the deal duration, and attempts to decode the payload CID and piece CID from the provided PieceCommitment object. If there is an error in the decoding process, the content status is updated to indicate that the proposal has failed and the error is returned.
- The code then creates a label for the deal proposal, and uses the Filecoin client to make the deal proposal. If there is any error during this process, the content status is updated to indicate that the proposal has failed and the error is returned.
- Assuming the proposal is successful, the code creates a new ContentDeal object and stores it in the database. It also updates the content status to indicate that the proposal is being sent. The proposal is then sent using a separate function, and any errors that occur during this process are handled and the content status is updated accordingly.
- If the proposal is successful, the code updates the content status to indicate that the proposal has been sent. If this is an end-to-end (E2E) content, the data transfer is initiated.