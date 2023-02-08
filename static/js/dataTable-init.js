function initProductTable()
{
    productTable = $('#productTable').DataTable({
        language: {
            url: "static/vendor/dataTables/languages/ch-tw.json"  
        },
        oLanguage: {
            "sProcessing": "<i class='fa fa-spinner fa-spin fa-3x fa-fw'></i><span class='sr-only'>Loading...</span>"
        },
        processing: true,
        serverSide:true,
        ajax: {
            url: "/api/products/getList",
            type: "GET",
        },
        columns:[
            {
                data: 'name', 
            },
            {
                data: null,
                render: function(data, type, row) {
                    output = '<div style="display:flex; align-items:center;">';
                    // input & updateBtn
                    output += '<input type="number" class="amountInput form-control" style="width:100px;" data-id="'+row.id+'" value="'+row.amount+'">';
                    output += '<button type="button" class="btn btn-success btn-sm mdfBtn_'+row.id+'" style="margin-left:5px" onclick="updateAmount('+row.id+')">修改</button>';
                    // loading icon
                    output += '<i class="fa fa-spinner fa-spin loading_'+row.id+'" style="margin-left:5px;font-size:20px;display:none;"></i>';
                    // success
                    output += '<i class="fa fa-solid fa-check statusIcon success_'+row.id+'" style="color:#12af55;margin-left:5px;font-size:20px;display:none;"></i>';
                    // failed
                    output += '<i class="fa fa-solid fa-xmark statusIcon failed_'+row.id+'" style="color:#FF0000;margin-left:5px;font-size:20px;display:none;"></i>';
                    output += '</div>';
                    output += '<div class="msgItem msgItem_'+row.id+'" style="display:none;">message: <span class="msg_'+row.id+'" style="color:#FF0000;"></span></div>';
                    return output;
                }
            },
            {
                data: 'amountNotice', 
                render: function(data, type, row) {
                    if(data>0) {
                        return data;
                    } else {
                        return "無";
                    }
                }
            },
            {
                data: 'formatTime',
            },
            {
                data: null, 
                render: function(data, type, row) {
                    output = '<div style="display:flex; align-items:center;">';
                    output += '<button type="button" class="btn btn-primary btn-sm" style="margin-right:5px" onclick="modifyProduct('+row.id+', \''+row.name+'\', '+row.amount+', '+row.amountNotice+')">編輯產品</button>';
                    output += '<button type="button" class="btn btn-danger btn-sm" onclick="doDelete(this, '+row.id+')">刪除產品</button>';
                    output += '</div>';
                    return output;
                }
            },
        ],
        columnDefs:[
            {
                targets: 4,
                orderable: false,
                responsivePriority: 2,
            },
            {
                targets: 0,
                responsivePriority: 1,
            }
        ],
    });
    return productTable;
}